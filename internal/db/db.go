package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/joho/godotenv"
)

func Connect() (*sql.DB, error) {
	if err := godotenv.Load("environment"); err != nil {
		return nil, err
	}

	dbURL := os.Getenv("DB_CONNECTION")
	if dbURL == "" {
		return nil, errors.New("DB_CONNECTION not found in environment file")
	}

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if err := applyMigrations(db); err != nil {
		return nil, err
	}

	//return &Database{DB: db}, nil
	return db, nil
}

func applyMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	switch {
	case err == nil:
		log.Println("migrations applied successfully")

		return nil
	case errors.Is(err, migrate.ErrNoChange):
		return nil
	default:
		return fmt.Errorf("migration failed: %w", err)
	}
}
