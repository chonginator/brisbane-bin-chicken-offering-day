package api

import (
	"net/http"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/resource"
)

type SuburbsPageData struct {
	Suburbs []resource.Resource
	Query   string
}

func (cfg *Config) HandlerSuburbs(w http.ResponseWriter, r *http.Request) {
	suburbs := cfg.suburbs

	query := r.URL.Query().Get("q")
	if r.URL.Query().Has("q") {
		suburbs = resource.FilterByName(cfg.suburbs, query)
	}

	cfg.respondWithHTML(w, "index.html", SuburbsPageData{
		Suburbs: suburbs,
		Query:   query,
	})
}
