package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"question/config"
	"question/handler"

	_ "github.com/mattn/go-sqlite3"
)

// internal function to create the gateway table
func createGatewayTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS gateway(
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
		Name TEXT type UNIQUE,
		IpAddress TEXT NOT NULL
	);
	`
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

// internal function to create the route table
func createRouteTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS route(
		Id INTEGER NOT NULL PRIMARY KEY,
		Prefix TEXT,
		GatewayId INTEGER
	);
	`
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func main() {

	// Performing basic database operations
	db := config.InitDB(config.Dbpath)
	defer db.Close()
	createGatewayTable(db)
	createRouteTable(db)

	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/gateway", handler.CreateGateway)
	http.HandleFunc("/route", handler.CreateRoute)
	http.HandleFunc("/search/route", handler.SearchRoute)
	http.ListenAndServe(":8000", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Replace this file with your own code and best of luck")
}
