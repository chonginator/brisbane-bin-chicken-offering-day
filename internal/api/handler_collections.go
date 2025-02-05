package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Collection struct {
	Day string
}

type CollectionsPageData struct {
	Collections []Collection
}

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

	collections := make([]Collection, len(dbCollections))
	for i, collection := range dbCollections {
		collections[i] = Collection{
			Day: collection.CollectionDay,
		}
	}

	data := CollectionsPageData{
		Collections: collections,
	}

	respondWithHTML(w, http.StatusOK, cfg.templates["collections.html"], data)
}