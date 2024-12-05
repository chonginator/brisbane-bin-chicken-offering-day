package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	// TODO: Import correct libsql client for Go
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	const suburbsDataFilePath = "data/suburbs-and-adjoining-suburbs.json"

	suburbs, err := loadAndProcessSuburbs(suburbsDataFilePath)
	if err != nil {
		log.Fatalf("Error loading suburbs data: %v", err)
	}

	fmt.Println(suburbs)

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

	err = seedSuburbs(db, suburbs)

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

func seedSuburbs(db *sql.DB, suburbNames []string) error {
	return nil
}