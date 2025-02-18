package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/chonginator/brisbane-bin-chicken-offering-day/sql/schema"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	goose.SetDialect("sqlite3")
	err = goose.Up(db, "sql/schema")
	if err != nil {
		log.Fatalf("Error running up migration: %v", err)
	}
}
