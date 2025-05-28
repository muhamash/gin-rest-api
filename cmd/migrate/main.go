package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Provide migration direction as an argument: up or down")
	}

	direction := os.Args[1]
	db, err :=  sql.Open("sqlite3", "./firstDatabase.db")

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to create SQLite instance: %v", err)
	}

	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations") 
	if err != nil {
		log.Fatalf("Failed to open migration files: %v", err)
	}
	
	m, err := migrate.NewWithInstance(
		"file",
		fSrc,
		"sqlite3",
		instance,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to revert migrations: %v", err)
		}
	default:
		log.Fatalf("Invalid migration direction: %s. Use 'up' or 'down'.", direction)		
	}	
}