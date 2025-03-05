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

	if r.Header.Get("HX-Request") == "true" {
		filteredSuburbs := filterSuburbs(cfg.suburbs, query)
		cfg.respondWithHTML(w, "suburbs-list", SuburbsPageData{
			Suburbs: filteredSuburbs,
			Query: query,
		})
		return
	}

	if r.URL.Query().Has("q") {
		filteredSuburbs := filterSuburbs(cfg.suburbs, query)
		cfg.respondWithHTML(w, "suburbs", SuburbsPageData{
			Suburbs: filteredSuburbs,
			Query: query,
		})
		return
	}

	cfg.respondWithHTML(w, "suburbs", SuburbsPageData{Suburbs: cfg.suburbs})
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
