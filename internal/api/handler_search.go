package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/resource"
)

type SearchData struct {
	Addresses []resource.Resource
	Query string
}

func (cfg *Config) HandlerSearch(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query().Get("q")
	if queryString == "" {
		return
	}
	query := sql.NullString{
		String: queryString,
		Valid: true,
	}

	dbAddresses, err := cfg.db.SearchAddresses(context.Background(), database.SearchAddressesParams{
		Limit: 10,
		Query: query,
	})
	if err != nil {
		cfg.respondWithError(w, http.StatusInternalServerError, "failed to search addresses", err)
		return
	}

	if len(dbAddresses) == 0 {
		cfg.respondWithHTML(w, "no_results_addresses_list.html", SearchData{
			Query: query.String,
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