package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/google/uuid"
)

type CollectionWeek struct {
	WeekStarting string `json:"week_starting"`
	Zone         string `json:"zone"`
}

func seedCollectionWeeks(dbQueries *database.Queries, filepath string) error {
	data, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer data.Close()

	collectionWeeks := []CollectionWeek{}
	decoder := json.NewDecoder(data)
	err = decoder.Decode(&collectionWeeks)
	if err != nil {
		return err
	}

	for _, record := range collectionWeeks {
		weekStarting, err := time.Parse(time.DateOnly, record.WeekStarting)
		if err != nil {
			return fmt.Errorf("couldn't parse bin collection week starting: %w", err)
		}

		_, err = dbQueries.CreateCollectionWeek(context.Background(), database.CreateCollectionWeekParams{
			ID:            uuid.New(),
			WeekStartDate: weekStarting,
			Zone:          record.Zone,
		})

		if err != nil {
			return fmt.Errorf("error creating bin collection week: %w", err)
		}
	}

	return nil
}
