package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {

	g := gin.Default()

	// Define the versions of the routes for the application
	v1 := g.Group("/api/v1")
	{
		// v1.POST("/users", app.createUser)
		// v1.GET("/users/:id", app.getUser)
		// v1.PUT("/users/:id", app.updateUser)
		// v1.DELETE("/users/:id", app.deleteUser)
		v1.POST("/events", app.createEvent)
		v1.GET("/events", app.getAllEvent)
		v1.GET("/events/:id", app.getEvent)
		v1.PUT("/events/:id", app.updateEvent)
		v1.DELETE("/events/:id", app.deleteEvent)
		// v1.POST("/events/:event_id/attendees", app.addAttendee)
		// v1.GET("/events/:event_id/attendees", app.getAttendees)
		// v1.DELETE("/events/:event_id/attendees/:attendee_id", app.removeAttendee)
		// v1.POST("/auth/register", app.registerUser)
	}

	return g
}