package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/muhamash/go-first-rest-api/internal/database"
)


func RetrieveUserFromContext(c *gin.Context) *database.User {
	userAny, exists := c.Get("user")
	if !exists {
		return &database.User{}
	}
	
	user, ok := userAny.(*database.User) 
	if !ok {
		return &database.User{}
	}

	return user
}