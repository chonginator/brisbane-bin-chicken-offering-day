package api

import "net/http"

func (cfg *Config) HandlerRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/suburbs", http.StatusMovedPermanently)
}