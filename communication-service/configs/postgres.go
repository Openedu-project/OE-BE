package configs

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func ConnectDatabase() {
	// Example connection string
	dsn := fmt.Sprintf("host=localhost port=5432 user postgres dbname=openedu sslmode=disable")
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Postgres pin failed: %v", err)
	}

	log.Println("Postgres connected")
}
