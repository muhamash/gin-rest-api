package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhamash/go-first-rest-api/internal/database"
)

type EventHandler struct {
	Models database.Models
}

// create event
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),  "status":"error"})
		return
	}

	err := h.Models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event", "detail": err.Error(), "status":"error"})
		return
	} 
	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"createdEvent":       event,
	})
}

// get all events
func (h *EventHandler) GetAllEvent(c *gin.Context) {
	events, err := h.Models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to retrieve events",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"totalEvents":  len(events),
		"events":       events,
	})
}

// get single event by Id
func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InvalId event Id", "detail": err.Error()})
		return
	}
	 
	event, err := h.Models.Events.GET(id)
	fmt.Println("Event is:", event, id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status":"ok","event": event})
}

// update event by Id
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event database.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "detail": err.Error()})
		return
	}
	event.Id = id

	// Track what was updated
	updatedFields := gin.H{}

	if event.Name != nil {
		updatedFields["name"] = *event.Name
	}
	if event.Description != nil {
		updatedFields["description"] = *event.Description
	}
	if event.Date != nil {
		updatedFields["date"] = event.Date.Format(time.RFC3339)
	}
	if event.Location != nil {
		updatedFields["location"] = *event.Location
	}
	if event.OwnerId != nil {
		updatedFields["ownerId"] = *event.OwnerId
	}

	if len(updatedFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields provided to update"})
		return
	}

	if err := h.Models.Events.Update(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"updatedEvent": updatedFields,
	})
}


// delete event by Id
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InvalId event Id"})
		return
	}

	existingEvent, err := h.Models.Events.GET(id)
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	} 

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if err := h.Models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": fmt.Sprintf("Event with ID %d deleted successfully", id),
		"deletedEvent": gin.H{
			"id":          existingEvent.Id,
			"name":        existingEvent.Name,
			"description": existingEvent.Description,
			"date":        existingEvent.Date.Format(time.RFC3339),
			"location":    existingEvent.Location,
			"ownerId":     existingEvent.OwnerId,
		},
		"deletedAt": time.Now().Format(time.RFC3339),
		// "deletedBy": c.GetString("userID"), 
	})
}