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

type SafeUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}


// create user
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

// get user utility function
func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

// get user by ID
func (m *UserModel) Get(id int) (*User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	return m.getUser(query, id)
}

// get user by email
func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	return m.getUser(query, email)
}

// get all users
func (m *UserModel) GetAllUser()([]*SafeUser, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "SELECT * FROM users"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*SafeUser{}
	for rows.Next(){
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil{
			return nil, err
		}

		safeUser := &SafeUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		users = append(users, safeUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err 
	}

	return users, nil
}