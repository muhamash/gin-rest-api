package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *application) serve() error {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: app.routes(),	
		IdleTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 3 << 20, // 1 MB
	}

	log.Printf("Starting server on %s", server.Addr)
	log.Printf("Server started at %s", time.Now().Format(time.RFC3339))
	log.Printf("Server idle timeout: %s", server.IdleTimeout)
	log.Printf("Server read timeout: %s", server.ReadTimeout)
	log.Printf("Server write timeout: %s", server.WriteTimeout)
	log.Printf("Server max header bytes: %d", server.MaxHeaderBytes)
	log.Printf("Server is ready to accept requests running at port %d", app.port)

	return server.ListenAndServe()
}