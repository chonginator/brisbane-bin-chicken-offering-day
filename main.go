package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const port = "8080"
	const rootFilePath = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(rootFilePath)))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
