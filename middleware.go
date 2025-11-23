package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired reads Authorization: Bearer <token> and loads user from DB.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format: Bearer <token>"})
			return
		}
		token := parts[1]

		user := findUserByToken(token)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// store a copy of user in context
		c.Set("user", *user)
		c.Next()
	}
}
