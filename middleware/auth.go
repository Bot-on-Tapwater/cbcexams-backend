package middleware

import (
	"net/http"
	"strings"

	"github.com/bot-on-tapwater/cbcexams-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth is a middleware function for Gin that validates JSON Web Tokens (JWT)
// from the "Authorization" header of incoming HTTP requests. It ensures that
// the token is present, properly formatted, and valid. If the token is invalid
// or missing, the middleware aborts the request with a 401 Unauthorized status.
//
// The middleware extracts the "user_id" claim from the token and sets it in
// the Gin context for downstream handlers to use.
//
// Usage:
// Add this middleware to your Gin router to protect routes that require
// authentication.
//
// Example:
// router := gin.Default()
// router.Use(JWTAuth())
// router.GET("/protected", protectedHandler)
//
// Returns:
// - HTTP 401 Unauthorized if the "Authorization" header is missing or the token is invalid.
// - Proceeds to the next handler if the token is valid.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		token, err := utils.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
