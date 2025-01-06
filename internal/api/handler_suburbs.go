package api

import (
	"html/template"
	"net/http"
)

func (cfg *Config) HandlerSuburbs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't load template", err)
		return
	}

	respondWithHTML(w, http.StatusOK, tmpl, cfg.suburbNames)
}