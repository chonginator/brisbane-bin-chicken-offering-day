package api

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *Config) respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	w.Header().Set("Content-Type", "text/html")

	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	cfg.templates.ExecuteTemplate(w, "error.html", msg)
}

func (cfg *Config) respondWithHTML(w http.ResponseWriter, templateName string, data any) {
	err := cfg.templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		http.Error(w, "Coudln't execute template", http.StatusInternalServerError)
		return
	}
}
