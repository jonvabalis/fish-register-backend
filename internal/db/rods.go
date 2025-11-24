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

func InsertRod(ctx context.Context, db *sql.DB, rod core.Rod) error {
	_, err := squirrel.Insert("rods").
		SetMap(map[string]any{
			"uuid":           rod.UUID,
			"nickname":       rod.Nickname,
			"brand":          rod.Brand,
			"purchase_place": rod.PurchasePlace,
			"users_uuid":     rod.UserUUID,
		}).
		RunWith(db).
		ExecContext(ctx)

	return err
}

func GetUserRods(ctx context.Context, db *sql.DB, uuid uuid.UUID) ([]core.Rod, error) {
	query := squirrel.Select("uuid", "nickname", "brand", "purchase_place", "users_uuid").
		From("rods").
		Where(squirrel.Eq{"users_uuid": uuid})

	rows, err := query.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rods := make([]core.Rod, 0)

	for rows.Next() {
		var r core.Rod

		if err := rows.Scan(
			&r.UUID,
			&r.Nickname,
			&r.Brand,
			&r.PurchasePlace,
			&r.UserUUID,
		); err != nil {
			return []core.Rod{}, fmt.Errorf("scanning: %w", err)
		}

		rods = append(rods, r)
	}

	if err := rows.Err(); err != nil {
		return []core.Rod{}, fmt.Errorf("iterating rows: %w", err)
	}

	return rods, nil
}

func GetRod(ctx context.Context, db *sql.DB, rUUID uuid.UUID) (core.Rod, error) {
	query := squirrel.
		Select("uuid", "nickname", "brand", "purchase_place", "users_uuid").
		From("rods").
		Where(squirrel.Eq{"uuid": rUUID})

	row := query.RunWith(db).QueryRowContext(ctx)

	var r core.Rod
	err := row.Scan(
		&r.UUID,
		&r.Nickname,
		&r.Brand,
		&r.PurchasePlace,
		&r.UserUUID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Rod{}, nil
		}
		return core.Rod{}, fmt.Errorf("scanning: %w", err)
	}

	return r, nil
}

func UpdateRod(ctx context.Context, db *sql.DB, rod core.Rod) error {
	_, err := squirrel.Update("rods").
		SetMap(map[string]any{
			"nickname":       rod.Nickname,
			"brand":          rod.Brand,
			"purchase_place": rod.PurchasePlace,
		}).
		Where(squirrel.Eq{"uuid": rod.UUID}).
		RunWith(db).
		ExecContext(ctx)

	return err
}

func DeleteRod(ctx context.Context, db *sql.DB, rUUID uuid.UUID) error {
	_, err := squirrel.
		Delete("rods").
		Where(squirrel.Eq{
			"uuid": rUUID,
		}).
		RunWith(db).
		ExecContext(ctx)

	return err
}
