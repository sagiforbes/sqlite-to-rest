# sqlite-to-rest
Simple application for exporting sqlite's database as rest (GET) call

#Building

Build the project by typing: `go build -o sqlite-to-rest`

#Running
```
sqlite-to-rest -f <database file name> -p <optional port, default: 4080>
```

#Query table

To query a table use this template endpoint:
```
http://<host ip>:<host port (default 4080)>/<table name>?start=<index>&length=<limit record count>
```

For example, say your server runs locally (localhost) and you want to fetch all records in __cars__ table:
```
curl "localhost:4080/cats"
```

Thats it. The results are in the following JSON format:
```
{
    "Columns": ["Array","of column names"],
    "ColumnTypes":["Data base column type"],
    "DataRows":[
        [value,value....],
        [value,value....],
    ]

}

```

Note that the DataRows is array of arrays of rows. Each row 


#### Limit respons

To limit the result of the query use the optional: `start` and `length` query parameter strings. Using our current example, if you want to skip over the first 100 records and return 20 records:
```
curl "http://localhost:4080/cars?start=100&length=20"
```
