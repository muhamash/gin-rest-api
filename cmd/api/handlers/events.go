package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhamash/go-first-rest-api/cmd/api/utils"
	"github.com/muhamash/go-first-rest-api/internal/database"
)

type EventHandler struct {
	Models database.Models
}

// create event

// CreateEvent creates a new event
//
//	@Summary		Creates a new event
//	@Description	Creates a new event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		database.Event	true	"Event"
//	@Success		201		{object}	database.Event
//	@Router			/api/v1/events [post]
//	@Security		BearerAuth

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

// GetEvents returns all events
//
//	@Summary		Returns all events
//	@Description	Returns all events
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]database.Event
//	@Router			/api/v1/events [get]

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

// GetEvent returns a single event
//
//	@Summary		Returns a single event
//	@Description	Returns a single event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [get]

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

// UpdateEvent updates an existing event
//
//	@Summary		Updates an existing event
//	@Description	Updates an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			event	body		database.Event	true	"Event"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [put]
//	@Security		BearerAuth

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Fetch the existing event from DB
	existingEvent, err := h.Models.Events.GET(id)
	if err != nil || existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Get current user from context
	contextUser := utils.RetrieveUserFromContext(c)

	// Check if the requester is the owner
	if existingEvent.OwnerId == nil || *existingEvent.OwnerId != contextUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of the event"})
		return
	}

	var updateData database.Event
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "detail": err.Error()})
		return
	}

	// Track and apply changes
	updatedFields := gin.H{}
	updateData.Id = id

	if updateData.Name != nil {
		existingEvent.Name = updateData.Name
		updatedFields["name"] = *updateData.Name
	}
	if updateData.Description != nil {
		existingEvent.Description = updateData.Description
		updatedFields["description"] = *updateData.Description
	}
	if updateData.Date != nil {
		existingEvent.Date = updateData.Date
		updatedFields["date"] = updateData.Date.Format(time.RFC3339)
	}
	if updateData.Location != nil {
		existingEvent.Location = updateData.Location
		updatedFields["location"] = *updateData.Location
	}

	if len(updatedFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields provided to update"})
		return
	}

	if err := h.Models.Events.Update(existingEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"updatedEvent": updatedFields,
	})
}



// delete event by Id

// DeleteEvent deletes an existing event
//
//	@Summary		Deletes an existing event
//	@Description	Deletes an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		204
//	@Router			/api/v1/events/{id} [delete]
//	@Security		BearerAuth

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