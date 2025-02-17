package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Collection struct {
	Day string
	BinWeekType BinWeekType
}

type CollectionsPageData struct {
	Collections []Collection
}

type BinWeekType string

const (
	YellowRecycling BinWeekType = "Yellow Recycling"
	GreenWaste BinWeekType = "Green Waste"
)

func (cfg *Config) HandlerCollections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	propertyID, ok := vars["property_id"]
	if !ok {
		err := fmt.Errorf("property_id is required")
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	dbCollections, err := cfg.db.GetCollectionSchedulesByPropertyID(context.Background(), propertyID)
	if err != nil {
		err := fmt.Errorf("couldn't get addresses by property_id: %w", err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch collections", err)
	}

	zoneForCurrentWeekRow, err := cfg.db.GetZoneForCurrentWeek(context.Background())
	if err != nil {
		err := fmt.Errorf("couldn't get zone for the current week: %w", err)
		respondWithError(w, http.StatusInternalServerError, "failed to get zone for the current week", err)
	}

	var binWeekType BinWeekType
	if zoneForCurrentWeekRow.Zone == dbCollections[0].Zone {
		binWeekType = YellowRecycling
	} else {
		binWeekType = GreenWaste
	}

	collections := make([]Collection, len(dbCollections))
	for i, collection := range dbCollections {
		collections[i] = Collection{
			Day: collection.CollectionDay,
			BinWeekType: binWeekType,
		}
	}

	data := CollectionsPageData{
		Collections: collections,
	}

	respondWithHTML(w, http.StatusOK, cfg.templates["collections.html"], data)
}
