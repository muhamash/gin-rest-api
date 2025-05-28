package main

import (
	"database/sql"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"

	"github.com/muhamash/go-first-rest-api/internal/database"
	"github.com/muhamash/go-first-rest-api/internal/env"
)

type application struct {
	port int
	models database.Models
	jwtSecret string
}

func main() {
	db, err := sql.Open("sqlite3", "./firstDatabase.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	
	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		models:    models,
		jwtSecret: env.GetEnvString("JWT_SECRET", "muhamash_secret"),
	}

	if err := app.serve(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	} 
}