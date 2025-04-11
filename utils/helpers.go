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
// ParseUserIDFromString parses a user ID string into a UUID object.
// It takes a Gin context and the user ID string as input parameters.
// If the parsing is successful, it returns the parsed UUID and a boolean value of true.
// If the parsing fails, it sends a JSON response with an HTTP 400 status code and an error message,
// and returns an empty UUID along with a boolean value of false.
//
// Parameters:
//   - c: The Gin context used to send an error response if parsing fails.
//   - userIDStr: The user ID string to be parsed into a UUID.
//
// Returns:
//   - uuid.UUID: The parsed UUID if successful, or an empty UUID if parsing fails.
//   - bool: A boolean indicating whether the parsing was successful (true) or not (false).
func ParseUserIDFromString(c *gin.Context, userIDStr string) (uuid.UUID, bool) {
	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return uuid.UUID{}, false
	}
	return parsedUserID, true

}
