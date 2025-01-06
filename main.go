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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable is not set")
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
	mux.HandleFunc("/", apiCfg.HandlerRoot)
	mux.HandleFunc("/suburbs", apiCfg.HandlerSuburbs)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
