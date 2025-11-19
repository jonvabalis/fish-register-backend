package db

import (
	"context"
	"database/sql"
	"errors"
	"fish-register-backend/internal/core"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

func InsertSpecies(ctx context.Context, db *sql.DB, species core.Species) error {
	_, err := squirrel.Insert("species").
		SetMap(map[string]any{
			"uuid":        species.UUID,
			"name":        species.Name,
			"description": species.Description,
		}).
		RunWith(db).
		ExecContext(ctx)

	return err
}

func GetMultipleSpecies(ctx context.Context, db *sql.DB) ([]core.Species, error) {
	query := squirrel.Select("uuid", "name", "description").From("species")

	rows, err := query.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	species := make([]core.Species, 0)

	for rows.Next() {
		var s core.Species

		if err := rows.Scan(
			&s.UUID,
			&s.Name,
			&s.Description,
		); err != nil {
			return nil, fmt.Errorf("scanning: %w", err)
		}

		species = append(species, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return species, nil
}

func GetSpecies(ctx context.Context, db *sql.DB, sUUID uuid.UUID) (core.Species, error) {
	query := squirrel.
		Select("uuid", "name", "description").
		From("species").
		Where(squirrel.Eq{"uuid": sUUID})

	row := query.RunWith(db).QueryRowContext(ctx)

	var s core.Species
	err := row.Scan(
		&s.UUID,
		&s.Name,
		&s.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Species{}, nil
		}
		return core.Species{}, fmt.Errorf("scanning: %w", err)
	}

	return s, nil
}

func UpdateSpecies(ctx context.Context, db *sql.DB, species core.Species) error {
	_, err := squirrel.Update("species").
		SetMap(map[string]any{
			"name":        species.Name,
			"description": species.Description,
		}).
		Where(squirrel.Eq{"uuid": species.UUID}).
		RunWith(db).
		ExecContext(ctx)

	return err
}
