package database

import (
	"context"
	"database/sql"
	"time"
)


type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id          int    `json:"id"`
	OwnerName        string `json:"ownerName" binding:"required, min=3,max=50"`
	Description string `json:"description" binding:"required, min=3,max=200"`
	Date        string `json:"date"`
	Location    string `json:"location" binding:"required, min=3,max=100"`
	OwnerId int    `json:"ownerId" binding:"required"`
}

// craete a new event
func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO events (owner_name, description, date, location, owner_id)
			  VALUES ($1, $2, $3, $4, $5)`
	// _, err := m.DB.ExecContext(ctx, query, event.OwnerName, event.Description, event.Date, event.Location, event.OwnerId)
	
	return m.DB.QueryRowContext(ctx, query, event.OwnerName, event.Description, event.Date, event.Location, event.OwnerId).Scan(&event.Id)
	
}

// get all events
func (m *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM events`
	
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.Id, &event.OwnerName, &event.Description, &event.Date, &event.Location, &event.OwnerId); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err 
	}

	return events, nil
}

// get single event by Id
func (m *EventModel) GET(Id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM events WHERE id = $1`
	
	var event Event
	err := m.DB.QueryRowContext(ctx, query, Id).Scan(&event.Id, &event.OwnerName, &event.Description, &event.Date, &event.Location, &event.OwnerId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}

	return &event, nil
}

// update event by Id
func (m *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

 	query := `UPDATE events SET owner_name = $1, description = $2, date = $3, location = $4, owner_id = $5 WHERE id = $6`

	_, err := m.DB.ExecContext(ctx, query, event.OwnerName, event.Description, event.Date, event.Location, event.OwnerId, event.Id)
	if err != nil {	
		return err
	}

	return nil
}

// delete event by Id
func (m *EventModel) Delete(Id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM events WHERE id = $1`

	_, err := m.DB.ExecContext(ctx, query, Id)
	if err != nil {
		return err
	}

	return nil
}

