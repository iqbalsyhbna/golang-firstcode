// internal/config/database.go
package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() *sql.DB {
	var err error
	
	// Konfigurasi koneksi MySQL
	connStr := "root:@tcp(localhost:3306)/golang_db"
	DB, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}

	log.Println("Successfully connected to MySQL database")
	return DB
}

func GetDB() *sql.DB {
	return DB
}