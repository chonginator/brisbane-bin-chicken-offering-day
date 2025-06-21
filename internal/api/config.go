package api

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/chonginator/brisbane-bin-chicken-offering-day/internal/database"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	if err != nil {
		return nil, fmt.Errorf("error getting suburbs from database: %w", err)
	}

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

func toSlugFromName(name string) string {
	return strings.Join(strings.Split(strings.ToLower(name), " "), "-")
}

func toNameFromSlug(slug string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.Join(strings.Split(slug, "-"), " "))
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
