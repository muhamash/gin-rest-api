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
		v1.POST("/events", app.event.CreateEvent)
		v1.GET("/events", app.event.GetAllEvent)
		v1.GET("/events/:id", app.event.GetEvent)
		v1.PUT("/events/:id", app.event.UpdateEvent)
		v1.DELETE("/events/:id", app.event.DeleteEvent)

		v1.POST("/auth/register", app.auth.RegisterUser)
	}

	return g
}