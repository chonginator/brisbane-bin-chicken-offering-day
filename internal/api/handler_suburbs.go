package api

import (
	"net/http"
)

type SuburbsPageData struct {
	SuburbNames []string
}

func (cfg *Config) HandlerSuburbs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	data := SuburbsPageData{
		SuburbNames: cfg.suburbNames,
	}

	respondWithHTML(w, http.StatusOK, cfg.templates["index.html"], data)
}
