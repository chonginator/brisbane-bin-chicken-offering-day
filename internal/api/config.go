package api

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
)

type Config struct {
	db        *database.Queries
	templates *template.Template
}

func NewAPIConfig(dbURL string) (*Config, error) {
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	dbQueries := database.New(db)

	templates, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	apiCfg := &Config{
		db:        dbQueries,
		templates: templates,
	}

	return apiCfg, nil
}

func parseTemplates() (*template.Template, error) {
	var templateFilepaths []string
	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			templateFilepaths = append(templateFilepaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing component templates: %w", err)
	}

	fmt.Println("Parsing template files:")
	for _, templateFilepath := range templateFilepaths {
		fmt.Println(templateFilepath)
	}

	tmpl, err := template.ParseFiles(templateFilepaths...)
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	return tmpl, nil
}
