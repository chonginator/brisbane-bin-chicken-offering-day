package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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

			log.Printf("Current address to insert: %+v", address)
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

func createAddressMap(dbQueries *database.Queries) (map[string]Address, error) {
	addressMap := make(map[string]Address)
	offset := 0

	for {
		dbAddresses, err := dbQueries.GetAddressBatch(context.Background(), database.GetAddressBatchParams{
			BatchSize: int64(defaultBatchSize),
			Offset:    int64(offset),
		})
		if err != nil {
			return nil, fmt.Errorf("couldn't get address batch from database: %w", err)
		}

		if len(dbAddresses) == 0 {
			break
		}

		for _, dbAddress := range dbAddresses {
			var unitNumber, houseNumberSuffix *string
			if dbAddress.UnitNumber.Valid {
				val := dbAddress.UnitNumber.String
				unitNumber = &val
			}
			if dbAddress.HouseNumberSuffix.Valid {
				val := dbAddress.HouseNumberSuffix.String
				houseNumberSuffix = &val
			}

			address := Address{
				ID:                dbAddress.ID,
				PropertyID:        dbAddress.PropertyID,
				UnitNumber:        unitNumber,
				HouseNumber:       dbAddress.HouseNumber,
				HouseNumberSuffix: houseNumberSuffix,
				StreetID:          dbAddress.StreetID,
				CollectionDay:     dbAddress.CollectionDay,
				Zone:              dbAddress.Zone,
			}

			key, err := fromAddressToKey(address)
			if err != nil {
				return nil, fmt.Errorf("error generating address key: %w", err)
			}

			addressMap[key] = address
		}

		offset += defaultBatchSize
	}

	log.Printf("Loaded %d existing addresses from database", len(addressMap))
	return addressMap, nil
}

func fromAddressToKey(address Address) (string, error) {
	var b strings.Builder

	unitNumber, houseNumberSuffix := "", ""
	if address.UnitNumber != nil {
		unitNumber = *address.UnitNumber
	}
	if address.HouseNumberSuffix != nil {
		houseNumberSuffix = *address.HouseNumberSuffix
	}

	writes := []struct {
		condition bool
		str       string
	}{
		{true, address.PropertyID},
		{address.UnitNumber != nil, unitNumber},
		{true, address.HouseNumber},
		{address.HouseNumberSuffix != nil, houseNumberSuffix},
		{true, address.StreetID.String()},
		{true, address.CollectionDay},
		{true, address.Zone},
	}
	sep := '_'

	for i, w := range writes {
		if w.condition {
			_, err := b.WriteString(w.str)
			if err != nil {
				return "", fmt.Errorf("failed to write string %s to address key: %w", w.str, err)
			}

			if i != len(writes)-1 {
				_, err := b.WriteRune(sep)
				if err != nil {
					return "", fmt.Errorf("failed to write separator %v to address key: %w", sep, err)
				}
			}
		}
	}

	return b.String(), nil
}

func createAddressesFromCollectionData(
	collectionRecords []CollectionRecord,
	suburbMap map[string]uuid.UUID,
	streetMap map[string]Street,
	addressMap map[string]Address,
) ([]Address, error) {
	var addresses []Address
	caser := cases.Title(language.English)

	for _, record := range collectionRecords {
		titleCasedStreetName := caser.String(record.StreetName)
		titleCasedSuburb := caser.String(record.Suburb)

		suburbID, ok := suburbMap[titleCasedSuburb]
		if !ok {
			return nil, fmt.Errorf("suburb not found: %s", titleCasedSuburb)
		}

		var streetID uuid.UUID
		streetKey := titleCasedStreetName + ":" + titleCasedSuburb

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
		addressKey, err := fromAddressToKey(address)
		if err != nil {
			return nil, err
		}

		if _, ok := addressMap[addressKey]; !ok {
			addresses = append(addresses, address)
		} else {
			log.Printf("Skipping existing address with key: %s", addressKey)
		}
	}

	return addresses, nil
}
