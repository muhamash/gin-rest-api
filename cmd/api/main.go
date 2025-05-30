package main

import (
	"database/sql"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"

	"github.com/muhamash/go-first-rest-api/cmd/api/handlers"
	"github.com/muhamash/go-first-rest-api/cmd/api/middleware"
	"github.com/muhamash/go-first-rest-api/internal/database"
	"github.com/muhamash/go-first-rest-api/internal/env"
)

type application struct {
	port int
	models database.Models
	jwtSecret string
	auth *handlers.AuthHandler
	event *handlers.EventHandler
	attendee *handlers.AttendeeHandler
	authMiddleware *middleware.AuthMiddleware
	// utils *utils.RetrieveUserFromContext
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
		auth:      &handlers.AuthHandler{Models: models},
		event: 	   &handlers.EventHandler{Models: models},
		attendee:  &handlers.AttendeeHandler{Models: models},
		authMiddleware:  &middleware.AuthMiddleware{Models: models},

		// utils : &ut
	}

	if err := app.serve(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	} 
}