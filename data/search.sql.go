// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: search.sql

package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const rankedSearch = `-- name: RankedSearch :many
SELECT uuid, display_name, country, locality, postal_code, street_address, region, external_id, active, created_at,
       ts_rank(search, websearch_to_tsquery('simple', $1)) rank
FROM account_groups
WHERE search @@ websearch_to_tsquery('simple', $1)
ORDER BY rank DESC
    LIMIT 25
`

type RankedSearchRow struct {
	Uuid          uuid.UUID      `json:"uuid"`
	DisplayName   string         `json:"display_name"`
	Country       string         `json:"country"`
	Locality      string         `json:"locality"`
	PostalCode    sql.NullString `json:"postal_code"`
	StreetAddress string         `json:"street_address"`
	Region        string         `json:"region"`
	ExternalID    string         `json:"external_id"`
	Active        bool           `json:"active"`
	CreatedAt     time.Time      `json:"created_at"`
	Rank          float32        `json:"rank"`
}

func (q *Queries) RankedSearch(ctx context.Context, websearchToTsquery string) ([]RankedSearchRow, error) {
	rows, err := q.db.QueryContext(ctx, rankedSearch, websearchToTsquery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RankedSearchRow
	for rows.Next() {
		var i RankedSearchRow
		if err := rows.Scan(
			&i.Uuid,
			&i.DisplayName,
			&i.Country,
			&i.Locality,
			&i.PostalCode,
			&i.StreetAddress,
			&i.Region,
			&i.ExternalID,
			&i.Active,
			&i.CreatedAt,
			&i.Rank,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const search = `-- name: Search :many
SELECT uuid, display_name, country, locality, postal_code, street_address, region, external_id, active, created_at
FROM account_groups
WHERE search @@ websearch_to_tsquery('simple', $1)
LIMIT 25
`

type SearchRow struct {
	Uuid          uuid.UUID      `json:"uuid"`
	DisplayName   string         `json:"display_name"`
	Country       string         `json:"country"`
	Locality      string         `json:"locality"`
	PostalCode    sql.NullString `json:"postal_code"`
	StreetAddress string         `json:"street_address"`
	Region        string         `json:"region"`
	ExternalID    string         `json:"external_id"`
	Active        bool           `json:"active"`
	CreatedAt     time.Time      `json:"created_at"`
}

func (q *Queries) Search(ctx context.Context, websearchToTsquery string) ([]SearchRow, error) {
	rows, err := q.db.QueryContext(ctx, search, websearchToTsquery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchRow
	for rows.Next() {
		var i SearchRow
		if err := rows.Scan(
			&i.Uuid,
			&i.DisplayName,
			&i.Country,
			&i.Locality,
			&i.PostalCode,
			&i.StreetAddress,
			&i.Region,
			&i.ExternalID,
			&i.Active,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}