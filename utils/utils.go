package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func DbCheckFile(dbFile string) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func DbExecCommand(dbFile string, sqlCommand string, parameters ...interface{}) (sql.Result, error) {
	var err error
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var cmd *sql.Stmt

	var sqlResult sql.Result

	cmd, err = db.Prepare(sqlCommand)
	if err != nil {
		return nil, err
	} else {
		sqlResult, err = cmd.Exec(parameters...)
		if err != nil {
			return nil, err
		}
	}
	return sqlResult, nil
}

type QueryResult struct {
	Columns     []string
	ColumnTypes []string
	DataRows    [][]interface{}
}

func DbQuery(dbFile string, query string, params ...interface{}) (*QueryResult, error) {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	ret := &QueryResult{
		Columns:     make([]string, len(columnTypes)),
		ColumnTypes: make([]string, len(columnTypes)),
		DataRows:    make([][]interface{}, 0),
	}

	for i, v := range columnTypes {
		ret.Columns[i] = v.Name()
		ret.ColumnTypes[i] = v.DatabaseTypeName()
	}

	var recData []interface{}
	var receiver []interface{}
	for rows.Next() {
		recData = make([]interface{}, len(columnTypes))
		receiver = make([]interface{}, len(columnTypes))
		for i := range columnTypes {
			receiver[i] = &recData[i]
		}
		err = rows.Scan(receiver...)
		if err != nil {
			return nil, err
		}

		ret.DataRows = append(ret.DataRows, recData)
	}
	return ret, nil
}

func DbCount(dbFile string, tableName string) (int64, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var ret int64
	rows.Scan(&ret)
	return ret, nil
}
