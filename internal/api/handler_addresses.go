package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Address struct {
	AddressString string
	Slug          string
}

type AddressesPageData struct {
	Addresses []Address
}

func (cfg *Config) HandlerAddresses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	streetSlug, ok := vars["street"]
	if !ok {
		err := fmt.Errorf("street parameter is required")
		cfg.respondWithError(w, http.StatusInternalServerError, err.Error(), err)
	}
	streetName := fromSlug(streetSlug)

	dbAddresses, err := cfg.db.GetAddressesByStreetName(context.Background(), streetName)
	if err != nil {
		err = fmt.Errorf("couldn't find addresses for %s: %w", streetName, err)
		cfg.respondWithError(w, http.StatusInternalServerError, "failed to fetch addresses", err)
		return
	}

	addresses := make([]Address, len(dbAddresses))

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

		addresses[i] = Address{
			Slug:          address.PropertyID,
			AddressString: addressString,
		}
	}

	if r.Header.Get("HX-Request") != "true" {
		cfg.respondWithHTML(w, "addresses.html", AddressesPageData{Addresses: addresses})
		return
	}

	if r.URL.Query().Has("q") {
		query := r.URL.Query().Get("q")
		filteredAddresses := filterAddresses(addresses, query)
		cfg.respondWithHTML(w, "addresses-list", AddressesPageData{Addresses: filteredAddresses})
		return
	}

	cfg.respondWithHTML(w, "addresses-partial", AddressesPageData{Addresses: addresses})
}

func filterAddresses(addresses []Address, query string) []Address {
	filtered := make([]Address, 0)
	for _, address := range addresses {
		if strings.Contains(strings.ToLower(address.AddressString), strings.ToLower(query)) {
			filtered = append(filtered, address)
		}
	}

	return filtered
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
