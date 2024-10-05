package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	// Konfigurasi koneksi MySQL
	connStr := "root:@tcp(localhost:3306)/golang_db"
	DB, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}
}
