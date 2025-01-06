package main

import (
	"context"
	"html/template"
	"net/http"
)

func (cfg *apiConfig) handlerSuburbs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	suburbs, err := cfg.db.GetSuburbs(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't load suburbs", err)
		return
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't load template", err)
		return
	}

	data := struct{
		SuburbNames []string
	}{}

	for _, suburb := range suburbs {
		data.SuburbNames = append(data.SuburbNames, suburb.Name)
	}

	respondWithHTML(w, http.StatusOK, tmpl, data)
}