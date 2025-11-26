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

func CreateCatch(ctx context.Context, db *sql.DB, catch core.Catch) error {
	insertMap := map[string]any{
		"uuid":       catch.UUID,
		"nickname":   catch.Nickname,
		"comment":    catch.Comment,
		"caught_at":  catch.CaughtAt,
		"created_at": catch.CreatedAt,
		"users_uuid": catch.UsersUUID,
	}

	// Handle nullable fields
	if catch.Length != nil {
		insertMap["length"] = *catch.Length
	} else {
		insertMap["length"] = nil
	}

	if catch.Weight != nil {
		insertMap["weight"] = *catch.Weight
	} else {
		insertMap["weight"] = nil
	}

	if catch.SpeciesUUID != nil {
		insertMap["species_uuid"] = *catch.SpeciesUUID
	} else {
		insertMap["species_uuid"] = nil
	}

	if catch.LocationsUUID != nil {
		insertMap["locations_uuid"] = *catch.LocationsUUID
	} else {
		insertMap["locations_uuid"] = nil
	}

	if catch.RodsUUID != nil {
		insertMap["rods_uuid"] = *catch.RodsUUID
	} else {
		insertMap["rods_uuid"] = nil
	}

	_, err := squirrel.Insert("catches").
		SetMap(insertMap).
		RunWith(db).
		ExecContext(ctx)

	return err
}

func GetCatch(ctx context.Context, db *sql.DB, catchUUID uuid.UUID) (core.Catch, error) {
	query := squirrel.
		Select(
			"uuid",
			"nickname",
			"length",
			"weight",
			"comment",
			"caught_at",
			"created_at",
			"species_uuid",
			"locations_uuid",
			"users_uuid",
			"rods_uuid",
		).
		From("catches").
		Where(squirrel.Eq{"uuid": catchUUID})

	row := query.RunWith(db).QueryRowContext(ctx)

	var c core.Catch
	err := row.Scan(
		&c.UUID,
		&c.Nickname,
		&c.Length,
		&c.Weight,
		&c.Comment,
		&c.CaughtAt,
		&c.CreatedAt,
		&c.SpeciesUUID,
		&c.LocationsUUID,
		&c.UsersUUID,
		&c.RodsUUID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Catch{}, nil
		}
		return core.Catch{}, fmt.Errorf("scanning: %w", err)
	}

	return c, nil
}

func GetUserCatches(ctx context.Context, db *sql.DB, userUUID uuid.UUID) ([]core.Catch, error) {
	query := squirrel.
		Select(
			"uuid",
			"nickname",
			"length",
			"weight",
			"comment",
			"caught_at",
			"created_at",
			"species_uuid",
			"locations_uuid",
			"users_uuid",
			"rods_uuid",
		).
		From("catches").
		Where(squirrel.Eq{"users_uuid": userUUID}).
		OrderBy("caught_at DESC")

	rows, err := query.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying catches: %w", err)
	}
	defer rows.Close()

	catches := make([]core.Catch, 0)

	for rows.Next() {
		var c core.Catch

		if err := rows.Scan(
			&c.UUID,
			&c.Nickname,
			&c.Length,
			&c.Weight,
			&c.Comment,
			&c.CaughtAt,
			&c.CreatedAt,
			&c.SpeciesUUID,
			&c.LocationsUUID,
			&c.UsersUUID,
			&c.RodsUUID,
		); err != nil {
			return []core.Catch{}, fmt.Errorf("scanning: %w", err)
		}

		catches = append(catches, c)
	}

	if err := rows.Err(); err != nil {
		return []core.Catch{}, fmt.Errorf("iterating rows: %w", err)
	}

	return catches, nil
}

func UpdateCatch(ctx context.Context, db *sql.DB, catchUUID uuid.UUID, updateData core.UpdateCatchData) error {
	updateMap := make(map[string]any)

	if updateData.Nickname != nil {
		updateMap["nickname"] = *updateData.Nickname
	}

	if updateData.Length != nil {
		updateMap["length"] = *updateData.Length
	}

	if updateData.Weight != nil {
		updateMap["weight"] = *updateData.Weight
	}

	if updateData.Comment != nil {
		updateMap["comment"] = *updateData.Comment
	}

	if updateData.CaughtAt != nil {
		updateMap["caught_at"] = *updateData.CaughtAt
	}

	if updateData.SpeciesUUID != nil {
		updateMap["species_uuid"] = *updateData.SpeciesUUID
	}

	if updateData.LocationsUUID != nil {
		updateMap["locations_uuid"] = *updateData.LocationsUUID
	}

	if updateData.RodsUUID != nil {
		updateMap["rods_uuid"] = *updateData.RodsUUID
	}

	if len(updateMap) == 0 {
		return nil
	}

	_, err := squirrel.Update("catches").
		SetMap(updateMap).
		Where(squirrel.Eq{"uuid": catchUUID}).
		RunWith(db).
		ExecContext(ctx)

	return err
}

func DeleteCatch(ctx context.Context, db *sql.DB, catchUUID uuid.UUID) error {
	_, err := squirrel.
		Delete("catches").
		Where(squirrel.Eq{"uuid": catchUUID}).
		RunWith(db).
		ExecContext(ctx)

	return err
}
