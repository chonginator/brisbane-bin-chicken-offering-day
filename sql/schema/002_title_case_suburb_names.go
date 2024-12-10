package schema

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
  goose.AddMigrationContext(UpTitleCase, DownTitleCase)
}

func UpTitleCase(ctx context.Context, tx *sql.Tx) error {
  rows, err := tx.QueryContext(ctx, `
    SELECT id, name FROM suburbs
  `)
  if err != nil {
    return err
  }
  defer rows.Close()

  caser := cases.Title(language.English)
  for rows.Next() {
    var id uuid.UUID
    var name string

    err := rows.Scan(&id, &name)
    if err != nil {
      return err
    }

    titleCasedName := caser.String(name)

    _, err = tx.ExecContext(ctx, `
      UPDATE suburbs
      SET name = ?
      WHERE id = ?
    `, titleCasedName, id)
    if err != nil {
      return err
    }
  }

  if rows.Err() != nil {
    return rows.Err()
  }

  return nil
}

func DownTitleCase(ctx context.Context, tx *sql.Tx) error {
  rows, err := tx.QueryContext(ctx, `
    SELECT id, name FROM suburbs
  `)
  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    var id uuid.UUID
    var name string

    err = rows.Scan(&id, &name)
    if err != nil {
      return err
    }

    upperCaseName := strings.ToUpper(name)

    _, err = tx.ExecContext(ctx, `
      UPDATE suburbs
      SET name = ?
      WHERE id = ?
    `, upperCaseName, id)

    if err != nil {
      return err
    }
  }

  if rows.Err() != nil {
    return rows.Err()
  }

  return nil
}