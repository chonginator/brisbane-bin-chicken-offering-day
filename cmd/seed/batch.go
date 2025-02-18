package main

import (
	"log"
	"time"
)

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
