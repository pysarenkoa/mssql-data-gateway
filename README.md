# mssql-data-gateway

mssql data access over http

**Compilation:**  
```go build -ldflags "-s -w"```

**Usage:**
```
 usage: cisco_data_gateway.exe <command>
        where <command> is one of
        install, remove, debug, start, stop, pause or continue.
```

HTTP POST [host.example.domain.com]:9090/sql_data  
REQUEST BODY:
```
{
    "credentials": {
        "username": "some_username",
        "password": "some_password",
        "database": "my_db",
        "host": "host.name",
        "port": 1433
    },
    "sql_query": "SELECT PersonID, EnterpriseName FROM my_db.dbo.Agent AS Agents WHERE EnterpriseName IN ('leliukh.v','akhtyrskyi.m');"
}
```

RESPONSE BODY:
```
[
    {
        "EnterpriseName": "leliukh.v",
        "PersonID": 5002
    },
    {
        "EnterpriseName": "akhtyrskyi.m",
        "PersonID": 5729
    }
]
```