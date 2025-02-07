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

type BinCollectionWeek struct {
	WeekStarting string `json:"week_starting"`
	Zone         string `json:"zone"`
}

func seedBinCollectionWeeks(dbQueries *database.Queries, filepath string) error {
	data, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer data.Close()

	binCollectionWeeks := []BinCollectionWeek{}
	decoder := json.NewDecoder(data)
	err = decoder.Decode(&binCollectionWeeks)
	if err != nil {
		return err
	}

	for _, record := range binCollectionWeeks {
		weekStarting, err := time.Parse(time.DateOnly, record.WeekStarting)
		if err != nil {
			return fmt.Errorf("couldn't parse bin collection week starting: %w", err)
		}

		_, err = dbQueries.CreateBinCollectionWeek(context.Background(), database.CreateBinCollectionWeekParams{
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
