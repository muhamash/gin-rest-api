package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhamash/go-first-rest-api/cmd/api/utils"
	"github.com/muhamash/go-first-rest-api/internal/database"
)

type AttendeeHandler struct {
	Models database.Models
}

// AddAttendeeToEvent adds an attendee to an event
// @Summary		Adds an attendee to an event
// @Description	Adds an attendee to an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		201		{object}	database.Attendee
// @Router			/api/v1/events/{id}/attendees/{userId} [post]
// @Security		BearerAuth

func (h *AttendeeHandler) RegisterAttendeeToEvent(c *gin.Context)  {
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID", "detail": err.Error()})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "detail": err.Error()})
		return
	}

	// fmt.Print(userId, eventId)

	event, err := h.Models.Events.GET(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event", "detail": err.Error()})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := h.Models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user", "detail": err.Error()})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	contextUser := utils.RetrieveUserFromContext(c)

	if event.OwnerId != nil && *event.OwnerId == contextUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot register yourself as an attendee for your own event"})
		return
	}	

	existingAttendee, err := h.Models.Attendees.GetByEventAndAttendee(eventId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing attendee", "detail": err.Error()})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already registered for this event"})
		return
	}

	attendee := database.Attendee{
		EventId: eventId,
		UserId:  userId,
	}
	
	result, err := h.Models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register attendee", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "User registered as an attendee for the event successfully",
		"attendee": result,
	})

}

// GetAttendeesForEvent returns all attendees for a given event
//
//	@Summary		Returns all attendees for a given event
//	@Description	Returns all attendees for a given event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	[]database.User
//	@Router			/api/v1/events/{id}/attendees [get]

func (h *AttendeeHandler) GetAttendeesForEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID", "detail": err.Error()})
		return
	}

	event, err := h.Models.Events.GET(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event", "detail": err.Error()})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Not found ant event",})
		return
	}


	attendees, err := h.Models.Attendees.GetAttendeesByEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"totalAttendees": len(attendees),
		"attendees":  attendees,
		"eventId": eventId,
		"eventName": event.Name,
		"eventLocation": event.Location,
		"ownerId": event.OwnerId,
	})
}

// GetEventsByAttendee returns all events for a given attendee
//
//	@Summary		Returns all events for a given attendee
//	@Description	Returns all events for a given attendee
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Attendee ID"
//	@Success		200	{object}	[]database.Event
//	@Router			/api/v1/attendees/{id}/events [get]

func (h *AttendeeHandler) GetEventsByAttendee(c *gin.Context) {
	user, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "detail": err.Error()})
		return
	}

	attendee, err := h.Models.Attendees.GetEventsByAttendee(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event", "detail": err.Error()})
		return
	}

	if attendee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"attendeeId": user,
		"event": attendee,
	})
}

// DeleteAttendeeFromEvent deletes an attendee from an event
// @Summary		Deletes an attendee from an event
// @Description	Deletes an attendee from an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		204
// @Router			/api/v1/events/{id}/attendees/{userId} [delete]
// @Security		BearerAuth

func (h *AttendeeHandler) DeleteAttendeeFromEvent(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "detail": err.Error()})
		return
	}

	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user eventId", "detail": err.Error()})
		return
	}

	event, err := h.Models.Events.GET(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// user := app.GetUserFromContext(c)
	// if event.OwnerId != user.Id {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "You are not authorized to delete an attendeeFromEvent"})
	// }

	err = h.Models.Attendees.Delete(userId, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "oka"})
}