package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Suburb struct {
	Name string `json:"suburb_name"`
}

func seedSuburbs(dbQueries *database.Queries, filepath string) error {
	if dbQueries == nil {
		return errors.New("database connection is nil")
	}

	suburbsDataFile, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer suburbsDataFile.Close()

	suburbsData := []Suburb{}
	decoder := json.NewDecoder(suburbsDataFile)
	err = decoder.Decode(&suburbsData)
	if err != nil {
		return err
	}

	suburbsSet := make(map[string]bool, len(suburbsData))
	for _, suburb := range suburbsData {
		suburbsSet[suburb.Name] = true
	}

	suburbMap, err := createSuburbMap(dbQueries)
	if err != nil {
		return fmt.Errorf("failed to create suburb map: %w", err)
	}

	caser := cases.Title(language.English)
	for name := range suburbsSet {
		titleCasedName := caser.String(name)
		if _, ok := suburbMap[titleCasedName]; ok {
			continue
		}

		_, err := dbQueries.CreateSuburb(context.Background(), database.CreateSuburbParams{
			ID:   uuid.New(),
			Name: titleCasedName,
		})

		if err != nil {
			return fmt.Errorf("failed to create suburb %s: %w", titleCasedName, err)
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
