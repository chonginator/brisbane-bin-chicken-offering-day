package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const port = "8080"

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
