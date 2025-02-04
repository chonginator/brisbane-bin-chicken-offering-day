package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func seedStreets(db *sql.DB, streets []Street) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := database.New(tx)

	startTime := time.Now()
	err = processBatch(streets, defaultBatchSize, func(batch []Street) error {
		for _, street := range batch {
			_, err := qtx.CreateStreet(context.Background(), database.CreateStreetParams{
				ID:       street.ID,
				Name:     street.Name,
				SuburbID: street.SuburbID,
			})
			if err != nil {
				return fmt.Errorf("error creating street: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	totalDuration := time.Since(startTime)
	log.Printf("Street seeding completed in: %v", totalDuration)

	return tx.Commit()
}

func createStreetMap(dbQueries *database.Queries) (map[string]Street, error) {
	caser := cases.Title(language.English)

	streetMap := make(map[string]Street)
	dbStreets, err := dbQueries.GetStreetsWithSuburb(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting streets from database: %w", err)
	}

	for _, street := range dbStreets {
		titleCasedStreetName := caser.String(street.StreetName)
		titleCasedSuburb := caser.String(street.SuburbName)

		streetKey := titleCasedStreetName + ":" + titleCasedSuburb

		streetMap[streetKey] = Street{
			ID:       street.ID,
			Name:     titleCasedStreetName,
			SuburbID: street.ID_2,
		}
	}

	return streetMap, nil
}