package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
ParseUserIDFromString parses a string into a UUID
and handles errors with a JSON response
*/
func ParseUserIDFromString(c *gin.Context, userIDStr string) (uuid.UUID, bool) {
	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return uuid.UUID{}, false
	}
	return parsedUserID, true
	
}