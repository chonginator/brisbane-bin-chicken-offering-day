package api

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
)

type Config struct {
	db        *database.Queries
	suburbs   []Suburb
	templates *template.Template
}

func NewAPIConfig(dbURL string) (*Config, error) {
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	dbQueries := database.New(db)

	dbSuburbs, err := dbQueries.GetSuburbs(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting suburbs from database: %w", err)
	}

	suburbs := make([]Suburb, 0, len(dbSuburbs))
	for _, suburb := range dbSuburbs {
		suburbs = append(suburbs, Suburb{
			Name: suburb.Name,
			Slug: toSlug(suburb.Name),
		})
	}

	sort.Slice(suburbs, func(i, j int) bool {
		return suburbs[i].Name < suburbs[j].Name
	})

	templates, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	apiCfg := &Config{
		db:        dbQueries,
		suburbs:   suburbs,
		templates: templates,
	}

	return apiCfg, nil
}

func toSlug(name string) string {
	return strings.Join(strings.Split(strings.ToLower(name), " "), "-")
}

func parseTemplates() (*template.Template, error) {
	pages, err := filepath.Glob("templates/pages/*.html")
	if err != nil {
		return nil, fmt.Errorf("error finding page templates: %w", err)
	}

	tmpl, err := template.ParseFiles(pages...)
	if err != nil {
		return nil, fmt.Errorf("error parsing page template: %w", err)
	}

	partials, err := filepath.Glob("templates/partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("error finding partial templates: %w", err)
	}

	tmpl, err = tmpl.ParseFiles(partials...)
	if err != nil {
		return nil, fmt.Errorf("error parsing partial templates: %w", err)
	}

	return tmpl, nil
}
