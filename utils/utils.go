package utils

import (
	"database/sql"

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

type QueryData []map[string]interface{}

func DbQuery(dbFile string, query string, params ...interface{}) (QueryData, error) {

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
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var ret QueryData
	var recData []interface{}
	var receiver []interface{}
	for rows.Next() {
		rec := make(map[string]interface{})
		recData = make([]interface{}, len(cols))
		receiver = make([]interface{}, len(cols))
		for i, _ := range cols {
			receiver[i] = &recData[i]
		}
		rows.Scan(receiver...)
		for i, col := range cols {
			rec[col] = recData[i]
		}
		ret = append(ret, rec)
	}
	return ret, nil
}
