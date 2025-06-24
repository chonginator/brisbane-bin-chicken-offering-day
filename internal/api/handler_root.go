package api

import "net/http"

func (cfg *Config) HandlerRoot(w http.ResponseWriter, r *http.Request) {
	cfg.respondWithHTML(w, "index.html", nil)
}
