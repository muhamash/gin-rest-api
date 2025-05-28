package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 1: Check for existing email
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
	if err := m.DB.QueryRowContext(ctx, checkQuery, user.Email).Scan(&exists); err != nil {
		return fmt.Errorf("email check failed: %w", err)
	}
	if exists {
		return fmt.Errorf("email already exists")
	}

	// Step 2: Insert the user
	insertQuery := `INSERT INTO users (username, password, email) VALUES (?, ?, ?)`
	result, err := m.DB.ExecContext(ctx, insertQuery, user.Username, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("LastInsertId failed: %w", err)
	}
	user.ID = int(id)

	// Step 3: Select the created user (if you want)
	// Optional: remove this if you already know the fields
	selectQuery := `SELECT username, email FROM users WHERE id = ?`
	err = m.DB.QueryRowContext(ctx, selectQuery, user.ID).Scan(&user.Username, &user.Email)
	if err != nil {
		return fmt.Errorf("failed to fetch created user: %w", err)
	}

	return nil
}

