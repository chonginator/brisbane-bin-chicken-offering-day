package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Street struct {
	Name string
	Slug string
}

type StreetsPageData struct {
	Streets []Street
}

func (cfg *Config) HandlerStreets(w http.ResponseWriter, r *http.Request) {
	suburbName := r.URL.Query().Get("suburbName")
	if suburbName == "" {
		err := fmt.Errorf("suburb parameter required")
		cfg.respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

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

	if r.Header.Get("HX-Request") != "true" {
		cfg.respondWithHTML(w, "streets.html", StreetsPageData{Streets: streets})
		return
	}

	if r.URL.Query().Has("q") {
		query := r.URL.Query().Get("q")
		filteredStreets := filterStreets(streets, query)
		cfg.respondWithHTML(w, "streets-list", StreetsPageData{Streets: filteredStreets})
		return
	}
	cfg.respondWithHTML(w, "streets-partial.html", StreetsPageData{Streets: streets})
}

func filterStreets(streets []Street, query string) []Street {
	filtered := make([]Street, 0)
	for _, street := range streets {
		if strings.Contains(strings.ToLower(street.Name), strings.ToLower(query)) {
			filtered = append(filtered, street)
		}
	}

	return filtered
}
