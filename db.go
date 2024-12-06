package main

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"log"

	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func runMigrations(db *sql.DB) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Path to migration files
		"sqlite3",           // Database type
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully.")
}

func initDB() *sql.DB {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return db
}
