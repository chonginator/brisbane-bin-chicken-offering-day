package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerRoot)
	mux.HandleFunc("/suburbs", handlerSuburbs)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/suburbs", http.StatusMovedPermanently)
}
