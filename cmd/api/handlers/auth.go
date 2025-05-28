package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhamash/go-first-rest-api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name    string `json:"name" binding:"required,min=3,max=50"`
}

type AuthHandler struct {
	Models database.Models
}

// registerHandler handles user registration requests.
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	fmt.Println("DB is:", h.Models.Users.DB)
	err := h.Models.Users.DB.Ping()
	if err != nil {
		log.Fatalf("DB Ping failed: %v", err)
	}


	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	req.Password = string(hashedPassword)
	
	user := database.User{
		Username: req.Name,
		Password: string(hashedPassword),
		Email:    req.Email,		
	}

	if err := h.Models.Users.Insert(&user); err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "detail": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
	
}	