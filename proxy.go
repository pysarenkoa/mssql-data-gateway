// +build windows

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"net/http"
	"strconv"
	"strings"
)

type RequestBody struct {
	SqlQuery    string `json:"sql_query"`
	Credentials struct {
		Database string `json:"database"`
		Username string `json:"username"`
		Password string `json:"Password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"credentials"`
}

func (r *RequestBody) checkRequest() error {
	if r.Credentials.Database == "" {
		return errors.New("request param [credentials.database] not found")
	}
	if r.Credentials.Username == "" {
		return errors.New("request param [credentials.username] not found")
	}
	if r.Credentials.Password == "" {
		return errors.New("request param [credentials.password] not found")
	}
	if r.Credentials.Host == "" {
		return errors.New("request param [credentials.host] not found")
	}
	if r.Credentials.Port == 0 {
		return errors.New("request param [credentials.port] not found")
	}
	if r.SqlQuery == "" {
		return errors.New("request param [sql_query] not found")
	}

	return nil
}

var DB *sql.DB
var req *RequestBody

func pingDbContext() error {
	ctx := context.Background()
	return DB.PingContext(ctx) // nil if connected
}

func initDbConnection() error {
	if DB != nil {
		return pingDbContext()
	}

	cr := req.Credentials
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		cr.Host, cr.Username, cr.Password, cr.Port, cr.Database)

	var err error
	DB, err = sql.Open("sqlserver", connStr)
	if err != nil {
		return err
	}

	return pingDbContext()
}

func writeDbDataToHttpResponse(w http.ResponseWriter) error {
	rows, err := DB.Query(req.SqlQuery)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	defer rows.Close()
	var res []map[string]interface{}

	count := len(columns)
	values := make([]sql.NullString, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
		rowData := make(map[string]interface{})
		for i, v := range values {
			// NOTE: FROM THE GO BLOG: JSON and GO - 25 Jan 2011:
			// The json package uses map[string]interface{} and []interface{} values to store arbitrary JSON objects and arrays;
			// it will happily unmarshal any valid JSON blob into a plain interface{} value. The default concrete Go types are:
			//
			// bool for JSON booleans,
			// float64 for JSON numbers,
			// string for JSON strings, and
			// nil for JSON null.
			if nx, ok := strconv.ParseFloat(v.String, 64); ok == nil {
				rowData[columns[i]] = nx
			} else if b, ok := strconv.ParseBool(v.String); ok == nil {
				rowData[columns[i]] = b
			} else if "string" == fmt.Sprintf("%T", v) {
				rowData[columns[i]] = v
			} else if "sql.NullString" == fmt.Sprintf("%T", v) {
				if v.String != "" {
					rowData[columns[i]] = v.String
				} else {
					rowData[columns[i]] = nil
				}
			} else {
				fmt.Printf("Failed on if for type %T of %v\n", v, v)
			}
		}

		res = append(res, rowData)
	}

	defer endOfWork()

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func runSqlOverHttpProxy() {
	// handle [/sql_data] route
	http.HandleFunc("/sql_data", func(w http.ResponseWriter, r *http.Request) {

		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-Awis-Session-Token, Authorization, X-Device-Type")
			return
		}

		// only POST
		if r.Method != "POST" {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		// only JSON
		var contentType = strings.ToLower(r.Header.Get("Content-Type"))
		if !strings.Contains(contentType, "application/json") {
			http.Error(w, "Header 'Content-Type: application/json' not found", http.StatusBadRequest)
			return
		}

		// decode request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// required params check
		err = req.checkRequest()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// db connection
		err = initDbConnection()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// read data and write response body
		err = writeDbDataToHttpResponse(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// starting http server
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		panic(err)
	}
}

// close connection and clean variables
func endOfWork() {
	_ = DB.Close()
	DB = nil
	req = nil
}
