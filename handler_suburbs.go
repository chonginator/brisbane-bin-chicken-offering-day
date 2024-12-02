package main

import (
	"html/template"
	"net/http"
)

func handlerSuburbs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't load template", err)
		return
	}

	data := struct {
		Items []string
	}{
		Items: []string{
			"Carindale",
			"Carina",
			"Carina Heights",
		},
	}

	respondWithHTML(w, http.StatusOK, tmpl, data)
}