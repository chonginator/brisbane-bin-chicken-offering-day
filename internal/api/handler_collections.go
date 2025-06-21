package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Collection struct {
	Date    string
	Message string
}

type CollectionsPageData struct {
	Collection Collection
}

var weekdayMap = map[string]time.Weekday{
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
	"Sunday":    time.Sunday,
}

func (cfg *Config) HandlerCollections(w http.ResponseWriter, r *http.Request) {
	propertyID := r.PathValue("property_id")
	if propertyID == "" {
		err := fmt.Errorf("property_id is required")
		cfg.respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	collectionSchedule, err := cfg.db.GetCollectionScheduleByPropertyID(context.Background(), propertyID)
	if err != nil {
		err := fmt.Errorf("couldn't get addresses by property_id: %w", err)
		cfg.respondWithError(w, http.StatusInternalServerError, "failed to fetch collections", err)
	}

	nextCollectionWeek, err := cfg.db.GetNextCollectionWeek(context.Background())
	if err != nil {
		err := fmt.Errorf("couldn't get the next collection week: %w", err)
		cfg.respondWithError(w, http.StatusInternalServerError, "failed to get next collection week", err)
	}

	dayDifference := weekdayMap[collectionSchedule.CollectionDay] - nextCollectionWeek.WeekStartDate.Weekday()
	nextCollectionDate := nextCollectionWeek.WeekStartDate.AddDate(0, 0, int(dayDifference))

	var message string
	if nextCollectionWeek.Zone == collectionSchedule.Zone {
		message = "General waste and yellow recycling bin day"
	} else {
		message = "General waste and green recycling bin day"
	}

	data := CollectionsPageData{
		Collection: Collection{
			Date:    nextCollectionDate.Format("Monday, 2 Jan 2006"),
			Message: message,
		},
	}

	cfg.respondWithHTML(w, "collections.html", data)
}
