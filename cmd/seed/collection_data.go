package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const defaultBatchSize = 50

type CollectionRecord struct {
	PropertyID        string  `json:"property_id"`
	UnitNumber        *string `json:"unit_number"`
	HouseNumber       string  `json:"house_number"`
	HouseNumberSuffix *string `json:"house_number_suffix"`
	StreetName        string  `json:"street_name"`
	Suburb            string  `json:"suburb"`
	CollectionDay     string  `json:"collection_day"`
	Zone              string  `json:"zone"`
}

type Street struct {
	ID       uuid.UUID
	Name     string
	SuburbID uuid.UUID
}

type Address struct {
	ID                uuid.UUID
	PropertyID        string
	UnitNumber        *string
	HouseNumber       string
	HouseNumberSuffix *string
	StreetID          uuid.UUID
	CollectionDay     string
	Zone              string
}

func seedCollectionData(db *sql.DB, filepath string) error {
	collectionData, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer collectionData.Close()

	collectionRecords := []CollectionRecord{}
	decoder := json.NewDecoder(collectionData)
	err = decoder.Decode(&collectionRecords)
	if err != nil {
		return err
	}

	dbQueries := database.New(db)

	suburbMap := make(map[string]uuid.UUID)
	dbSuburbs, err := dbQueries.GetSuburbs(context.Background())
	if err != nil {
		return fmt.Errorf("error getting suburbs: %w", err)
	}
	for _, suburb := range dbSuburbs {
		suburbMap[suburb.Name] = suburb.ID
	}

	caser := cases.Title(language.English)

	streetMap := make(map[string]Street)
	dbStreets, err := dbQueries.GetStreetsWithSuburb(context.Background())
	if err != nil {
		return fmt.Errorf("error getting streets: %w", err)
	}
	for _, street := range dbStreets {
		titleCasedStreetName := caser.String(street.StreetName)
		titleCasedSuburb := caser.String(street.SuburbName)

		streetKey := titleCasedStreetName + ":" + titleCasedSuburb

		streetMap[streetKey] = Street{
			ID: street.ID,
			Name: titleCasedStreetName,
			SuburbID: street.ID_2,
		}
	}

	var addresses []Address

	for _, record := range collectionRecords {
		titleCasedStreetName := caser.String(record.StreetName)
		titleCasedSuburb := caser.String(record.Suburb)

		streetKey := titleCasedStreetName + ":" + titleCasedSuburb

		suburbID, ok := suburbMap[titleCasedSuburb]
		if !ok {
			return fmt.Errorf("suburb not found: %s", titleCasedSuburb)
		}

		var streetID uuid.UUID
		if street, ok := streetMap[streetKey]; ok {
			streetID = street.ID
		} else {
			streetID = uuid.New()
			streetMap[streetKey] = Street{
				ID:       streetID,
				Name:     titleCasedStreetName,
				SuburbID: suburbID,
			}
		}

		address := Address{
			ID:                uuid.New(),
			PropertyID:        record.PropertyID,
			UnitNumber:        record.UnitNumber,
			HouseNumber:       record.HouseNumber,
			HouseNumberSuffix: record.HouseNumberSuffix,
			StreetID:          streetID,
			CollectionDay:     caser.String(record.CollectionDay),
			Zone:              caser.String(record.Zone),
		}

		addresses = append(addresses, address)
	}

	streets := make([]Street, 0, len(streetMap))
	for _, street := range streetMap {
		streets = append(streets, street)
	}

	// err = seedStreets(db, streets)
	// if err != nil {
	// 	return err
	// }

	err = seedAddresses(db, addresses)
	if err != nil {
		return err
	}

	return nil
}

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

func seedAddresses(db *sql.DB, addresses []Address) error {
	
	startTime := time.Now()
	err := processBatch(addresses, defaultBatchSize, func(batch []Address) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		qtx := database.New(tx)

		for _, address := range batch {
			var houseNumberSuffix sql.NullString
			if address.HouseNumberSuffix != nil {
				houseNumberSuffix = sql.NullString{
					String: *address.HouseNumberSuffix,
					Valid:  true,
				}
			}

			var unitNumber sql.NullString
			if address.UnitNumber != nil {
				unitNumber = sql.NullString{
					String: *address.UnitNumber,
					Valid:  true,
				}
			}

			_, err := qtx.CreateAddress(context.Background(), database.CreateAddressParams{
				ID:                address.ID,
				PropertyID:        address.PropertyID,
				UnitNumber:        unitNumber,
				HouseNumber:       address.HouseNumber,
				HouseNumberSuffix: houseNumberSuffix,
				StreetID:          address.StreetID,
				CollectionDay:     address.CollectionDay,
				Zone:              address.Zone,
			})
			if err != nil {
				log.Printf("Failed on propertyID: %s, street: %s", address.PropertyID, address.StreetID)
				return fmt.Errorf("error creating address: %w", err)
			}
		}

		return tx.Commit()
	})

	if err != nil {
		return err
	}

	totalDuration := time.Since(startTime)
	log.Printf("Address seeding completed in: %v", totalDuration)
	return nil
}

func processBatch[T any](items []T, batchSize int, process func([]T) error) error {
	log.Printf("Starting to process %d items...", len(items))
	for i := 0; i < len(items); i += batchSize {
		batchStartTime := time.Now()

		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		if err := process(items[i:end]); err != nil {
			return err
		}

		batchDuration := time.Since(batchStartTime)
		percentComplete := float64(end) / float64(len(items)) * 100

		log.Printf("Batch completed: %.2f%% (%d/%d items processed). Batch took: %v", percentComplete, end, len(items), batchDuration)
	}
	return nil
}
