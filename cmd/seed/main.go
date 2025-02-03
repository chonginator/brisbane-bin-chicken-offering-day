package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	suburbsPath := filepath.Join(dataDir, "suburbs-and-adjoining-suburbs.json")
	suburbs, err := loadAndProcessSuburbs(suburbsPath)
	if err != nil {
		log.Fatalf("Error loading suburbs data: %v", err)
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
	dbQueries := database.New(db)

	err = seedSuburbs(dbQueries, suburbs)
	if err != nil {
		log.Fatalf("Error seeding suburbs: %v", err)
	}

	collectionsPath := filepath.Join(dataDir, "waste-collection-days-collection-days.json")
	err = seedCollectionData(db, collectionsPath)
	if err != nil {
		log.Fatalf("Error loading collections data: %v", err)
	}
}
