package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	const suburbsDataFilePath = "data/suburbs-and-adjoining-suburbs.json"

	suburbs, err := loadAndProcessSuburbs(suburbsDataFilePath)
	if err != nil {
		log.Fatalf("Error loading suburbs data: %v", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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

}

type Suburb struct {
	SuburbName string `json:"suburb_name"`
}

func loadAndProcessSuburbs(filename string) ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	suburbsDataFullFilePath := filepath.Join(dir, filename)
	suburbsDataFile, err := os.ReadFile(suburbsDataFullFilePath)
	if err != nil {
		return nil, err
	}

	suburbsData := []Suburb{}
	err = json.Unmarshal(suburbsDataFile, &suburbsData)
	if err != nil {
		return nil, err
	}

	suburbsSet := make(map[string]struct{})

	for _, suburb := range suburbsData {
		suburbsSet[suburb.SuburbName] = struct{}{}
	}

	suburbs := []string{}

	for suburbName := range suburbsSet {
		suburbs = append(suburbs, suburbName)
	}

	return suburbs, nil
}

func seedSuburbs(db *database.Queries, suburbNames []string) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	for _, suburbName := range suburbNames {
			_, err := db.CreateSuburb(context.Background(), database.CreateSuburbParams{
				ID: uuid.New(),
				Name: suburbName,
			})

			if err != nil {
				return err
			}
	}

	return nil
}