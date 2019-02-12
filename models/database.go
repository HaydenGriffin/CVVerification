package models

import (
	"database/sql"
	"fmt"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	fmt.Printf("Connecting to MySQL database \n")
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
