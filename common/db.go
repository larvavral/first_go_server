package common

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	// Create connection pool
	var err error
	db, err = sql.Open("mssql", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}