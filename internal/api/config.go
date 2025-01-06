package api

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
)

type Config struct {
	db *database.Queries
	suburbNames []string
}

func NewAPIConfig(dbURL string) (*Config, error) {
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	suburbs, err := dbQueries.GetSuburbs(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting suburbs from database: %v", err)
	}

	suburbNames := make([]string, 0, len(suburbs))
	for _, suburb := range suburbs {
		suburbNames = append(suburbNames, suburb.Name)
	}

	apiCfg := &Config{
		db: dbQueries,
		suburbNames: suburbNames,
	}

	return apiCfg, nil
}

