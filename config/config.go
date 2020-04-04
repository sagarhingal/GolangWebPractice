package config

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const Dbpath = "tanla.db"

// InitDB: a config function for intialising the database
func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}
