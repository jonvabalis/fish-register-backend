package db

import (
	"context"
	"database/sql"
	"fish-register-backend/internal/core"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func InsertLocation(ctx context.Context, db *sql.DB, location core.Location) error {
	_, err := squirrel.Insert("locations").
		SetMap(map[string]any{
			"uuid":    location.UUID,
			"name":    location.Name,
			"address": location.Address,
			"type":    location.Type,
		}).
		RunWith(db).
		ExecContext(ctx)

	return err
}

func GetLocations(ctx context.Context, db *sql.DB) ([]core.Location, error) {
	query := squirrel.Select("uuid", "name", "address", "type").From("locations")

	rows, err := query.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]core.Location, 0)

	for rows.Next() {
		var l core.Location

		if err := rows.Scan(
			&l.UUID,
			&l.Name,
			&l.Address,
			&l.Type,
		); err != nil {
			return nil, fmt.Errorf("scanning: %w", err)
		}

		locations = append(locations, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return locations, nil
}
