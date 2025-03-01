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
	w.Header().Set("Content-Type", "text/html")

	data := SuburbsPageData{
		Suburbs: cfg.suburbs,
	}

	respondWithHTML(w, http.StatusOK, cfg.templates["index.html"], data)
}
