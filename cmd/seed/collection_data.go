package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/google/uuid"
)

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
	ID         uuid.UUID
	Name string
	SuburbID   uuid.UUID
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

func seedCollectionData(db *database.Queries, filepath string) error {
	collectionData, err := os.Open(filepath)
	if err != nil {
		return err
	}

	collectionRecords := []CollectionRecord{}
	decoder := json.NewDecoder(collectionData)
	err = decoder.Decode(&collectionRecords)
	if err != nil {
		return err
	}

	suburbMap := make(map[string]uuid.UUID)
	suburbs, err := db.GetSuburbs(context.Background())
	if err != nil {
		return fmt.Errorf("error getting suburbs: %w", err)
	}
	for _, suburb := range suburbs {
		suburbMap[suburb.Name] = suburb.ID
	}

	streetMap := make(map[string]Street)

	var addresses []Address
	for _, record := range collectionRecords {
		streetKey := record.StreetName + ":" + record.Suburb

		suburbID, ok := suburbMap[record.Suburb]
		if !ok {
			return fmt.Errorf("suburb not found: %s", record.Suburb)
		}

		var streetID uuid.UUID
		if street, ok := streetMap[streetKey]; ok {
			streetID = street.ID
		} else {
			streetID = uuid.New()
			streetMap[streetKey] = Street{
				ID:         streetID,
				Name: record.StreetName,
				SuburbID:   suburbID,
			}
		}

		address := Address{
			ID:                uuid.New(),
			PropertyID:        record.PropertyID,
			UnitNumber:        record.UnitNumber,
			HouseNumber:       record.HouseNumber,
			HouseNumberSuffix: record.HouseNumberSuffix,
			StreetID:          streetID,
			CollectionDay:     record.CollectionDay,
			Zone:              record.Zone,
		}

		addresses = append(addresses, address)
	}

	streets := make([]Street, 0, len(streetMap))
	for _, street := range streetMap {
		streets = append(streets, street)
	}

	err = seedStreets(db, streets)
	if err != nil {
		return err
	}

	err = seedAddresses(db, addresses)
	if err != nil {
		return err
	}

	return nil
}

func seedStreets(db *database.Queries, streets []Street) error {
	for _, street := range streets {
		_, err := db.CreateStreet(context.Background(), database.CreateStreetParams{
			ID: street.ID,
			Name: street.Name,
			SuburbID: street.SuburbID,
		})
		if err != nil {
			return fmt.Errorf("error creating street: %w", err)
		}
	}
	return nil
}

func seedAddresses(db *database.Queries, addresses []Address) error {
	return nil
}
