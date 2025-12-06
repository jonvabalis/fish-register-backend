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

func UpdateUser(ctx context.Context, db *sql.DB, user core.UserAuth) error {
	_, err := squirrel.Update("users").
		SetMap(map[string]any{
			"username": user.Username,
			"email":    user.Email,
			"password": user.Password,
		}).
		Where(squirrel.Eq{"uuid": user.UUID}).
		RunWith(db).
		ExecContext(ctx)

	return err
}

func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (core.UserAuth, error) {
	query := squirrel.
		Select("uuid", "username", "email", "password", "role").
		From("users").
		Where(squirrel.Eq{"email": email})

	row := query.RunWith(db).QueryRowContext(ctx)

	var ua core.UserAuth
	err := row.Scan(
		&ua.UUID,
		&ua.Username,
		&ua.Email,
		&ua.Password,
		&ua.Role,
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

func GetUser(ctx context.Context, db *sql.DB, uuid uuid.UUID) (core.UserAuth, error) {
	query := squirrel.
		Select("uuid", "username", "email", "password").
		From("users").
		Where(squirrel.Eq{"uuid": uuid})

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

func GetUsers(ctx context.Context, db *sql.DB) ([]core.User, error) {
	query := squirrel.
		Select("uuid", "username", "email").
		From("users")

	rows, err := query.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]core.User, 0)

	for rows.Next() {
		var user core.User

		if err := rows.Scan(
			&user.UUID,
			&user.Username,
			&user.Email,
		); err != nil {
			return []core.User{}, fmt.Errorf("scanning: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return []core.User{}, fmt.Errorf("iterating rows: %w", err)
	}

	return users, nil
}

func DeleteUser(ctx context.Context, db *sql.DB, userUUID uuid.UUID) error {
	_, err := squirrel.
		Delete("users").
		Where(squirrel.Eq{
			"uuid": userUUID,
		}).
		RunWith(db).
		ExecContext(ctx)

	return err
}
