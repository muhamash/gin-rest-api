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
		
		v1.GET("/events", app.event.GetAllEvent)
		v1.GET("/events/:id", app.event.GetEvent)


		v1.POST("/auth/register", app.auth.RegisterUser)
		v1.GET("/auth/users", app.auth.GetAllUsers)
		v1.POST("/auth/login", app.auth.LoginUser)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.authMiddleware.RequireAuth())
	{
		authGroup.POST("/events", app.event.CreateEvent)
		authGroup.PUT("/events/:id", app.event.UpdateEvent)
		authGroup.DELETE("/events/:id", app.event.DeleteEvent)
		authGroup.POST("/events/:eventId/attendees/:userId", app.attendee.RegisterAttendeeToEvent)
		authGroup.GET("/events/attendees/:eventId", app.attendee.GetAttendeesForEvent)
		authGroup.GET("/attendees/events/:userId", app.attendee.GetEventsByAttendee)
		authGroup.DELETE("/events/attendees/:eventId/:userId", app.attendee.DeleteAttendeeFromEvent)
	}

	return g
}