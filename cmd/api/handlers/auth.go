package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
	jwtSecret string
}

type loginRequest struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type loginResponse struct{
	Token string `json: "token"`
}

// login
func (h *AuthHandler) LoginUser(c *gin.Context){
	var auth loginRequest

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}	

	existingUser, err := h.Models.Users.GetByEmail(auth.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}	
	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email not valid", "details": err.Error()})
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login error", "details": err.Error()})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUser.ID,
		"email":   existingUser.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return   
	}

	c.JSON(http.StatusOK, gin.H{
		"token": loginResponse{Token: tokenString},
		"userId": existingUser.ID,
		"userName": existingUser.Username,
	})
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

// get all users
func (h *AuthHandler) GetAllUsers(c *gin.Context){
	users, err := h.Models.Users.GetAllUser()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to retrieve users",
			"detail": err.Error(),
		})
		return
	}

	fmt.Print(users)

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"totalUsers":  len(users),
		"users":       users,
	})

}
