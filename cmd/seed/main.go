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

	suburbsDataPath := filepath.Join(dataDir, "suburbs-and-adjoining-suburbs.json")
	err = seedSuburbs(dbQueries, suburbsDataPath)
	if err != nil {
		log.Fatalf("Error seeding suburbs data: %v", err)
	}

	binCollectionWeeksDataPath := filepath.Join(dataDir, "waste-collection-days-collection-weeks.json")
	err = seedBinCollectionWeeks(dbQueries, binCollectionWeeksDataPath)
	if err != nil {
		log.Fatalf("Error seeding bin collection weeks data: %v", err)
	}

	collectionsPath := filepath.Join(dataDir, "waste-collection-days-collection-days.json")
	err = seedCollectionData(db, collectionsPath)
	if err != nil {
		log.Fatalf("Error loading collections data: %v", err)
	}
}
