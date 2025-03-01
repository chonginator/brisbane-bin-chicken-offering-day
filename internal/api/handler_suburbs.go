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
}

func (cfg *Config) HandlerSuburbs(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		query := r.URL.Query().Get("q")
		filteredSuburbs := filterSuburbs(cfg.suburbs, query)
		cfg.respondWithHTML(w, "suburbs-list", SuburbsPageData{Suburbs: filteredSuburbs})
		return
	}

	cfg.respondWithHTML(w, "index.html", SuburbsPageData{Suburbs: cfg.suburbs})
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
