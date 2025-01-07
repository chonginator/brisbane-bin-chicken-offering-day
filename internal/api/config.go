package api

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
)

type Config struct {
	db          *database.Queries
	suburbNames []string
	templates   map[string]*template.Template
}

func NewAPIConfig(dbURL string) (*Config, error) {
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	suburbs, err := dbQueries.GetSuburbs(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting suburbs from database: %w", err)
	}

	suburbNames := make([]string, 0, len(suburbs))
	for _, suburb := range suburbs {
		suburbNames = append(suburbNames, suburb.Name)
	}

	templates, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	apiCfg := &Config{
		db:          dbQueries,
		suburbNames: suburbNames,
		templates:   templates,
	}

	return apiCfg, nil
}

func parseTemplates() (map[string]*template.Template, error) {
	files, err := filepath.Glob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("error finding templates: %w", err)
	}

	templates := make(map[string]*template.Template)
	layoutFile := "templates/layout.html"

	for _, file := range files {
		if file == layoutFile {
			continue
		}

		name := filepath.Base(file)
		tmpl, err := template.ParseFiles(layoutFile, file)
		if err != nil {
			return nil, fmt.Errorf("error parsing template: %w", err)
		}

		templates[name] = tmpl
	}

	fmt.Println(templates)
	return templates, nil
}
