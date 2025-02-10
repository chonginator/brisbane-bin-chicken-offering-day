package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/google/uuid"
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
	data, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer data.Close()

	collectionRecords := []CollectionRecord{}
	decoder := json.NewDecoder(data)
	err = decoder.Decode(&collectionRecords)
	if err != nil {
		return err
	}

	dbQueries := database.New(db)

	suburbMap, err := createSuburbMap(dbQueries)
	if err != nil {
		return fmt.Errorf("error creating suburb map: %w", err)
	}

	streetMap, err := createStreetMap(dbQueries)
	if err != nil {
		return fmt.Errorf("error creating street map: %w", err)
	}

	streets := make([]Street, 0, len(streetMap))
	for _, street := range streetMap {
		streets = append(streets, street)
	}

	err = seedStreets(db, streets)
	if err != nil {
		return err
	}

	addressMap, err := createAddressMap(dbQueries)
	if err != nil {
		return fmt.Errorf("error creating address map: %w", err)
	}

	log.Printf("Processing %d collection records", len(collectionRecords))

	addresses, err := createAddressesFromCollectionData(collectionRecords, suburbMap, streetMap, addressMap)
	if err != nil {
		return fmt.Errorf("error creating addresses from collection data: %w", err)
	}

	log.Printf("Found %d new addresses to seed", len(addresses))

	err = seedAddresses(db, addresses)
	if err != nil {
		return err
	}

	return nil
}
