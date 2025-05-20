package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/api"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL environment variable is not set")
	}

	apiCfg, err := api.NewAPIConfig(dbURL)
	if err != nil {
		log.Fatalf("Error initializing API config: %v", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", apiCfg.HandlerRoot)
	mux.HandleFunc("/suburbs", apiCfg.HandlerSuburbs)
	mux.HandleFunc("/suburbs/{suburb}/streets", apiCfg.HandlerStreets)
	mux.HandleFunc("/suburbs/{suburb}/streets/{street}/addresses", apiCfg.HandlerAddresses)
	mux.HandleFunc("/suburbs/{suburb}/streets/{street}/addresses/{property_id}/collections", apiCfg.HandlerCollections)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
