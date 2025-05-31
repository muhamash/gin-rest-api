package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	redisclient "github.com/muhamash/go-first-rest-api/internal"
	"github.com/muhamash/go-first-rest-api/internal/database"
	"github.com/redis/go-redis/v9"
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
	Redis     *redis.Client
}
type loginRequest struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type loginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName"`
}

// login

// Login logs in a user
//
//	@Summary		Logs in a user
//	@Description	Logs in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body	loginRequest	true	"User"
//	@Success		200	{object}	loginResponse
//	@Router			/api/v1/auth/login [post]

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var auth loginRequest

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := h.Models.Users.GetByEmail(auth.Email)
	if err != nil || existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUser.ID,
		"email":   existingUser.Email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	accessString, err := accessToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUser.ID,
		"email":   existingUser.Email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	refreshString, err := refreshToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	key := fmt.Sprintf("refresh_token:user:%d", existingUser.ID)
	err = h.Redis.Set(redisclient.Ctx, key, refreshString, 7*24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Token:        accessString,
		RefreshToken: refreshString,
		UserId:       int64(existingUser.ID),
		UserName:     existingUser.Username,
	})
}


// @Summary      Register a new user
// @Description  Register a new user and return a JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  models.UserRegisterRequest  true  "User Data"
// @Success      200   {object}  models.UserResponse
// @Failure      400   {object}  models.ErrorResponse
// @Router       /auth/register [post]

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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token", "details": err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	userId := int64(claims["user_id"].(float64))
	key := fmt.Sprintf("refresh_token:user:%d", userId)

	storedToken, err := h.Redis.Get(redisclient.Ctx, key).Result()
	if err != nil || storedToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token mismatch"})
		return
	}

	// Issue new access token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"email":   claims["email"],
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	newTokenString, err := newToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// save or update the refresh token
	err = h.Redis.Set(redisclient.Ctx, key, newTokenString, 7*24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to updating refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newTokenString})
}


// get all users
// @Summary		Get all users
// @Description	Get all users
// @Tags			test
// @Accept			json
// @Produce		json
// @Success		200	{object}	gin.H "Successful response"
// @Router			/api/v1/auth/users [get]
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
