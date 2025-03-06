package api

import (
	"net/http"
	"strings"
)

type Suburb struct {
	Name string
	Slug string
}

type SuburbsPageData struct {
	Suburbs []Suburb
	Query string
}

func (cfg *Config) HandlerSuburbs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	suburbs := cfg.suburbs

	if r.URL.Query().Has("q") {
		suburbs = filterSuburbs(cfg.suburbs, query)
	}

	cfg.respondWithHTML(w, "index.html", SuburbsPageData{
		Suburbs: suburbs,
		Query: query,
	})
}

func filterSuburbs(suburbs []Suburb, query string) []Suburb {
	filtered := make([]Suburb, 0)

	for _, suburb := range suburbs {
		if strings.Contains(strings.ToLower(suburb.Name), strings.ToLower(query)) {
			filtered = append(filtered, suburb)
		}
	}

	return filtered
}
