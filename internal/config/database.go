// internal/config/database.go
package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DBMap map[string]*sql.DB

func InitDBs() {
	DBMap = make(map[string]*sql.DB)

	// Konfigurasi untuk multiple databases
	dbConfigs := map[string]string{
		"golang_db": "root:@tcp(localhost:3306)/golang_db",
	}

	for dbName, connStr := range dbConfigs {
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			log.Fatalf("Error opening database %s: %q", dbName, err)
		}

		// Set connection pool settings
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		// Verify connection
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error connecting to the database %s: %q", dbName, err)
		}

		DBMap[dbName] = db
		log.Printf("Successfully connected to MySQL database: %s", dbName)
	}
}

func GetDB(dbName string) *sql.DB {
	return DBMap[dbName]
}
