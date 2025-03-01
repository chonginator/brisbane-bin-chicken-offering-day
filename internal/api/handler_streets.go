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

	suburbSlug, ok := vars["suburb"]
	if !ok {
		err := fmt.Errorf("suburb parameter required")
		cfg.respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	suburbName := fromSlug(suburbSlug)
	fmt.Println("Getting streets for suburb:", suburbName)

	dbStreets, err := cfg.db.GetStreetsBySuburbName(context.Background(), suburbName)
	if err != nil {
		err = fmt.Errorf("couldn't find streets for %s: %w", suburbName, err)
		cfg.respondWithError(w, http.StatusInternalServerError, "failed to fetch streets", err)
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

	cfg.respondWithHTML(w, "streets.html", data)
}
