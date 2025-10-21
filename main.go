package main

import (
	"database/sql"
	"errors"
	"fish-register-backend/internal/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("environment"); err != nil {
		log.Fatal("error loading environment file: ", err)
	}

	dbURL := os.Getenv("DB_CONNECTION")
	if dbURL == "" {
		log.Fatal("DB_CONNECTION not found in environment file: ")
	}

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("application started successfully")

	router := gin.Default()

	router.GET("/register", handlers.Register)

	err = router.Run(":1111")
	if err != nil {
		return
	}
}

func runMigrations(db *sql.DB) error {
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
		return err
	}
}
