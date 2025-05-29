package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)


type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id          int    `json:"id"`
	Name   		*string `json:"name" min:"3" max:"50"`
	Description *string `json:"description" min:"3" max:"200"`
	Date        *time.Time `json:"date"`
	Location    *string `json:"location" min:"3" max:"100"`
	OwnerId     *int    `json:"ownerId"`
}

// craete a new event
func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO events (name, description, date, location, owner_id)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	// _, err := m.DB.ExecContext(ctx, query, event.name, event.Description, event.Date, event.Location, event.OwnerId)
	
	return m.DB.QueryRowContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.OwnerId).Scan(&event.Id)
	
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
		if err := rows.Scan(&event.Id, &event.Name, &event.OwnerId, &event.Description, &event.Date, &event.Location); err != nil {
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
	err := m.DB.QueryRowContext(ctx, query, Id).Scan(
		&event.Id,
		&event.Name,
		&event.OwnerId,
		&event.Description,
		&event.Date,
		&event.Location,
	)
	
	fmt.Println("Query is:", query, "Id is:", Id, err, event)

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

	setClauses := []string{}
	args := []interface{}{}
	argID := 1

	if event.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argID))
		args = append(args, *event.Name)
		argID++
	}
	if event.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argID))
		args = append(args, *event.Description)
		argID++
	}
	if event.Date != nil {
		setClauses = append(setClauses, fmt.Sprintf("date = $%d", argID))
		args = append(args, *event.Date)
		argID++
	}
	if event.Location != nil {
		setClauses = append(setClauses, fmt.Sprintf("location = $%d", argID))
		args = append(args, *event.Location)
		argID++
	}
	if event.OwnerId != nil {
		setClauses = append(setClauses, fmt.Sprintf("owner_id = $%d", argID))
		args = append(args, *event.OwnerId)
		argID++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add final ID condition
	query := fmt.Sprintf(`UPDATE events SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), argID)
	args = append(args, event.Id)

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
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

