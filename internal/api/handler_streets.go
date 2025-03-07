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
	SuburbName string
	SuburbSlug string
	Query string
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

	streets := make([]Street, len(dbStreets))
	for i, street := range dbStreets {
		streets[i] = Street{
			Name: street.Name,
			Slug: toSlug(street.Name),
		}
	}

	query := r.URL.Query().Get("q")
	if r.URL.Query().Has("q") {
		streets = filterStreets(streets, query)
	}

	cfg.respondWithHTML(w, "streets.html", StreetsPageData{
		Streets: streets,
		Query: query,
		SuburbName: suburbName,
		SuburbSlug: suburbSlug,
	})
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
