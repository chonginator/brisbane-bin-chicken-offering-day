package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Street struct {
	Name string
	Slug string
}

type StreetsPageData struct {
	Streets []Street
}

func (cfg *Config) HandlerStreets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	suburbSlug := vars["suburb"]
	suburbName := fromSlug(suburbSlug)

	dbStreets, err := cfg.db.GetStreetsBySuburbName(context.Background(), suburbName)
	if err != nil {
		err = fmt.Errorf("couldn't find streets for %s: %w", suburbName, err)
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	streets := make([]Street, len(dbStreets))
	for i, street := range dbStreets {
		streets[i] = Street{
			Name: street.Name,
			Slug: toSlug(street.Name),
		}
	}

	data := StreetsPageData{
		Streets: streets,
	}

	respondWithHTML(w, http.StatusOK, cfg.templates["streets.html"], data)
}
