package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main(){
	databaseURL := os.Getenv("DATABASE_URL")
  if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	migrationPath := "file://app/migrations"

	m, err := migrate.New(
		migrationPath,
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	log.Println("Starting database migration...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply.")
	} else {
		log.Println("Migrations applied successfully.")
	}
	
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNoChange {
		log.Printf("Warning: Could not get migration version: %v", err)
	} else {
		log.Printf("Current database version: %d (dirty: %t)", version, dirty)
	}
}