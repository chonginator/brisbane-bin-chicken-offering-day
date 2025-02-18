package api

import (
	"html/template"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Coudln't load error template", http.StatusInternalServerError)
		return
	}
	respondWithHTML(w, code, tmpl, msg)
}

func respondWithHTML(w http.ResponseWriter, code int, tmpl *template.Template, data any) {
	w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(code)
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Coudln't execute template", http.StatusInternalServerError)
		return
	}
}
