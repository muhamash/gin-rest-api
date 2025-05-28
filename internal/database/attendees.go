package database

import "database/sql"

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID        int    `json:"id"`
	EventID   int    `json:"eventId" binding:"required"` 
	UserId  int    `json:"userId" binding:"required"` 
}	