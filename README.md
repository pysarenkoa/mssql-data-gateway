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
        "username": "DevCore_local",
        "password": "some_password",
        "database": "np_awdb",
        "host": "PCCE-HDS1A",
        "port": 1433
    },
    "sql_query": "SELECT PersonID, EnterpriseName FROM np_awdb.dbo.t_Agent AS Agents WHERE EnterpriseName IN ('leliukh.v','akhtyrskyi.m');"
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