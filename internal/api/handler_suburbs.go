package api

import (
	"html/template"
	"net/http"
)

type SuburbsPageData struct {
	SuburbNames []string
}

func (cfg *Config) HandlerSuburbs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// TODO: Move template parsing outside of handler?
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't load template", err)
		return
	}

	data := SuburbsPageData{
		SuburbNames: cfg.suburbNames,
	}

	respondWithHTML(w, http.StatusOK, tmpl, data)
}
