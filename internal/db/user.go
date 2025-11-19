package db

import (
	"context"
	"database/sql"
	"errors"
	"fish-register-backend/internal/core"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func CreateUser(ctx context.Context, db *sql.DB, user core.UserAuth) error {
	_, err := squirrel.Insert("users").
		SetMap(map[string]any{
			"uuid":     user.UUID,
			"username": user.Username,
			"email":    user.Email,
			"password": user.Password,
		}).
		RunWith(db).
		ExecContext(ctx)

	return err
}

func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (core.UserAuth, error) {
	query := squirrel.
		Select("uuid", "username", "email", "password").
		From("users").
		Where(squirrel.Eq{"email": email})

	row := query.RunWith(db).QueryRowContext(ctx)

	var ua core.UserAuth
	err := row.Scan(
		&ua.UUID,
		&ua.Username,
		&ua.Email,
		&ua.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.UserAuth{}, nil
		}
		return core.UserAuth{}, fmt.Errorf("scanning: %w", err)
	}

	return ua, nil
}

func GetUserByUsername(ctx context.Context, db *sql.DB, username string) (core.UserAuth, error) {
	query := squirrel.
		Select("uuid", "username", "email", "password").
		From("users").
		Where(squirrel.Eq{"username": username})

	row := query.RunWith(db).QueryRowContext(ctx)

	var ua core.UserAuth
	err := row.Scan(
		&ua.UUID,
		&ua.Username,
		&ua.Email,
		&ua.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.UserAuth{}, nil
		}
		return core.UserAuth{}, fmt.Errorf("scanning: %w", err)
	}

	return ua, nil
}
