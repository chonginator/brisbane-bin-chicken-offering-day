package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
		caser := cases.Title(language.English)
		suburbs = append(suburbs, caser.String(suburbName))
	}

	return suburbs, nil
}

func seedSuburbs(db *database.Queries, suburbNames []string) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	for _, suburbName := range suburbNames {
		_, err := db.CreateSuburb(context.Background(), database.CreateSuburbParams{
			ID:   uuid.New(),
			Name: suburbName,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func createSuburbMap(dbQueries *database.Queries) (map[string]uuid.UUID, error) {
	suburbMap := make(map[string]uuid.UUID)
	dbSuburbs, err := dbQueries.GetSuburbs(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting suburbs from database: %w", err)
	}
	for _, suburb := range dbSuburbs {
		suburbMap[suburb.Name] = suburb.ID
	}
	return suburbMap, nil
}