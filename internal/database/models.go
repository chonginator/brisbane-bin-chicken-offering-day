// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID                uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PropertyID        string
	UnitNumber        sql.NullString
	HouseNumber       string
	HouseNumberSuffix sql.NullString
	StreetID          uuid.UUID
	CollectionDay     string
	Zone              string
}

type AddressSearch struct {
	PropertyID string
	SearchText string
}

type BinCollectionWeek struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	WeekStartDate time.Time
	Zone          string
}

type Street struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	SuburbID  uuid.UUID
}

type Suburb struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}
