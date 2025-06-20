package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/resource"
)

type SearchData struct {
	Addresses []resource.Resource
	Query string
}

func (cfg *Config) HandlerAddresses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		cfg.respondWithHTML(w, "placeholder_results.html", nil)
		return
	}

	safeQuery := strings.ReplaceAll(query, `"`, `""`)
	ftsQuery := fmt.Sprintf("\"%s\"*", safeQuery)
	dbAddresses, err := cfg.db.SearchAddresses(context.Background(), database.SearchAddressesParams{
		Limit: 10,
		Query: ftsQuery,
	})
	if err != nil {
		cfg.respondWithError(w, http.StatusInternalServerError, "failed to search addresses", err)
		return
	}

	if len(dbAddresses) == 0 {
		cfg.respondWithHTML(w, "no_results.html", SearchData{
			Query: query,
		})
	}

	addresses := make([]resource.Resource, len(dbAddresses))

	for i, row := range dbAddresses {
		addresses[i] = resource.Resource{
			Slug: row.PropertyID,
			Name: row.FormattedAddress,
		}
	}

	cfg.respondWithHTML(w, "addresses_list.html", SearchData{
		Addresses: addresses,
	})
}