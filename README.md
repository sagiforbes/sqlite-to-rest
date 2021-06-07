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
```bash
curl "localhost:4080/cats"
```

Thats it. The results are in the following JSON format:
```json
{
    "Columns": ["Array","of column names"],
    "ColumnTypes":["Data base column type"],
    "Rows":[
        [value,value....],
        [value,value....],
        .
        .
        .
    ]
}

```

Note that the Rows is array of arrays of rows. Each row 

#### Getting response in json format
If you want to get the Rows as an array of json object, each object is a row in the database, just prepend the end point with json. for example:
```bash
curl "http://localhost:4080/json/cars
```
The output is the following JSON format:
```json
{
    "ColumnTypes":["Data base column type"],
    "Rows":[
        {
            "field name": <Field value>,
            "field name": <Field value>,
            .
            .
            .
        }
    ]
}

```

#### Getting response in CSV format
If you want to get the Rows as an array of json object, each object is a row in the database, just prepend the end point with json. for example:
```bash
curl "http://localhost:4080/json/cars
```
The output is the following CSV format:
```csv
<Field1 name>, <Field2 name>
<Row 1 Field1>, <Row1 Field2>
<Row 2 Field1>, <Row2 Field2>
.
.
.


```

#### Limit respons

To limit the result of the query use the optional: `start` and `length` query parameter strings. Using our current example, if you want to skip over the first 100 records and return 20 records:
```bash
curl "http://localhost:4080/cars?start=100&length=20"
```

