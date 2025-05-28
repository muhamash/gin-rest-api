package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhamash/go-first-rest-api/internal/database"
)

// create event
func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event", "detail": err.Error()})
		return
	} 
	c.JSON(http.StatusCreated, event)
}

// get all event
func (app *application) getAllEvent(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// get single event by Id
func (app *application) getEvent(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InvalId event Id"})
		return
	}
	 
	event, err := app.models.Events.GET(Id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// update event by Id
func (app *application) updateEvent(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InvalId event Id"})
		return
	}

	existingEvent, err := app.models.Events.GET(Id)
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	updateEvent := &database.Event{}
	if err := c.ShouldBindJSON(updateEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateEvent.Id = Id
	if err := app.models.Events.Update(updateEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return		
	}

	c.JSON(http.StatusOK, updateEvent)
}

// delete event by Id
func (app *application) deleteEvent(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InvalId event Id"})
		return
	}

	existingEvent, err := app.models.Events.GET(Id)
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	} 

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if err := app.models.Events.Delete(Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}