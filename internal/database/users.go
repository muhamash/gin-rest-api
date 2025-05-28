package database

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}