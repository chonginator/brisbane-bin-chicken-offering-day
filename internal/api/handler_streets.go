package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/resource"
)

type StreetsPageData struct {
	Streets    []resource.Resource
	SuburbName string
	SuburbSlug string
	Query      string
}

func (cfg *Config) HandlerStreets(w http.ResponseWriter, r *http.Request) {
	suburbName := r.URL.Query().Get("suburbName")
	if suburbName == "" {
		err := fmt.Errorf("suburb parameter required")
		cfg.respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	suburbSlug := toSlug(suburbName)

	dbStreets, err := cfg.db.GetStreetsBySuburbName(context.Background(), suburbName)
	if err != nil {
		err = fmt.Errorf("couldn't find streets for %s: %w", suburbName, err)
		cfg.respondWithError(w, http.StatusInternalServerError, "failed to fetch streets", err)
		return
	}

	streets := make([]resource.Resource, len(dbStreets))
	for i, street := range dbStreets {
		streets[i] = resource.Resource{
			Name: street.Name,
			Slug: toSlug(street.Name),
		}
	}

	query := r.URL.Query().Get("q")
	if r.URL.Query().Has("q") {
		streets = resource.FilterByName(streets, query)
	}

	cfg.respondWithHTML(w, "streets.html", StreetsPageData{
		Streets:    streets,
		Query:      query,
		SuburbName: suburbName,
		SuburbSlug: suburbSlug,
	})
}
