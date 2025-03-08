package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/resource"
)

type AddressesPageData struct {
	Addresses []resource.Resource
}

func (cfg *Config) HandlerAddresses(w http.ResponseWriter, r *http.Request) {
	streetName := r.URL.Query().Get("streetName")
	if streetName == "" {
		err := fmt.Errorf("street name parameter is required")
		cfg.respondWithError(w, http.StatusInternalServerError, err.Error(), err)
	}

	dbAddresses, err := cfg.db.GetAddressesByStreetName(context.Background(), streetName)
	if err != nil {
		err = fmt.Errorf("couldn't find addresses for %s: %w", streetName, err)
		cfg.respondWithError(w, http.StatusInternalServerError, "failed to fetch addresses", err)
		return
	}

	addresses := make([]resource.Resource, len(dbAddresses))

	for i, address := range dbAddresses {
		var unitNumber, houseNumberSuffix string
		if address.UnitNumber.Valid {
			unitNumber = address.UnitNumber.String
		}
		if address.HouseNumberSuffix.Valid {
			houseNumberSuffix = address.HouseNumberSuffix.String
		}

		addressString, err := toAddressString(unitNumber, address.HouseNumber, houseNumberSuffix, streetName)
		if err != nil {
			err = fmt.Errorf("couldn't build address string: %w", err)
			cfg.respondWithError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}

		addresses[i] = resource.Resource{
			Slug: address.PropertyID,
			Name: addressString,
		}
	}

	query := r.URL.Query().Get("q")
	if r.URL.Query().Has("q") {
		addresses = resource.FilterByName(addresses, query)
	}

	cfg.respondWithHTML(w, "addresses.html", AddressesPageData{
		Addresses: addresses,
	})
}

func toAddressString(unitNumber, houseNumber, houseNumberSuffix, streetName string) (string, error) {
	var b strings.Builder

	writes := []struct {
		condition bool
		val       string
	}{
		{unitNumber != "", unitNumber + "/"},
		{true, houseNumber},
		{houseNumberSuffix != "", houseNumberSuffix},
		{true, " " + streetName},
	}

	for _, w := range writes {
		if w.condition {
			if _, err := b.WriteString(w.val); err != nil {
				return "", fmt.Errorf("failed to write address: %w", err)
			}
		}
	}

	return b.String(), nil
}
