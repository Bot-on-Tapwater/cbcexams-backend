package controllers

import (
	"fmt"
	"net/http"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
	"github.com/bot-on-tapwater/cbcexams-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UsersController struct {
	DB *gorm.DB
}

// Profile handles the retrieval of a user's profile based on the user ID
// stored in the request context. It performs the following steps:
//  1. Retrieves the `user_id` from the request context.
//  2. Validates that the `user_id` is a string.
//  3. Parses the `user_id` into a UUID format.
//  4. Queries the database for a user with the corresponding UUID.
//  5. Returns the user's profile as a JSON response if found, or an error
//     message if the user is not found or the `user_id` is invalid.
//
// Parameters:
// - c (*gin.Context): The Gin context containing the request and response objects.
//
// Responses:
// - 200 OK: Returns the user's profile in JSON format.
// - 400 Bad Request: Returns an error message if the `user_id` is invalid.
// - 404 Not Found: Returns an error message if the user is not found in the database.
func (uc *UsersController) Profile(c *gin.Context) {
	/* Retrieve the user_id from the context */
	userID := c.MustGet("user_id")

	/* Assert that userID is a string */
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	fmt.Printf("UserID is: %s\n", userIDStr)

	/* Convert userID to uuid.UUID */
	parsedUserID, ok := utils.ParseUserIDFromString(c, userIDStr)

	if !ok {
		return
	}

	/* Query the database using the parsed UUID */
	var user models.User
	if err := uc.DB.First(&user, "id = ?", parsedUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfile handles the HTTP request to update a user's profile.
// It performs the following steps:
// 1. Extracts the user ID from the JWT middleware.
// 2. Validates and parses the input JSON payload into an UpdateProfileInput struct.
// 3. Converts the user ID string to a UUID format.
// 4. Fetches the user record from the database using the parsed UUID.
// 5. Applies the updates to the user's profile using the provided input.
// 6. Saves the updated user record back to the database, ensuring email uniqueness.
// 7. Returns a success response with the updated user details or an appropriate error response.
//
// @param c *gin.Context - The Gin context containing the HTTP request and response.
// @response 200 OK - Profile updated successfully with updated user details.
// @response 400 Bad Request - Invalid user ID format or input payload.
// @response 404 Not Found - User not found in the database.
// @response 409 Conflict - Email already exists in the database.
func (uc *UsersController) UpdateProfile(c *gin.Context) {
	/* Get user ID from JWT middleware */
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	/* Parse Input */
	var input models.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	/* Convert userID to uuid */
	parsedUserID, ok := utils.ParseUserIDFromString(c, userID)

	if !ok {
		return
	}

	/* Fetch user from DB */
	var user models.User
	if err := uc.DB.First(&user, parsedUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	/* Apply updates */
	if err := user.UpdateProfile(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* Save to DB */
	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email alread exists"}) /* Handle email uniqueness */
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		"user": gin.H{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
		},
	})
}

// GetUsers handles the retrieval of user data based on the provided query parameters.
// If a "user_id" query parameter is provided, it fetches and returns the specific user
// corresponding to that ID. If no "user_id" is provided, it fetches and returns all users.
//
// Query Parameters:
//   - user_id (optional): The UUID of the user to retrieve.
//
// Responses:
//   - 200 OK: Returns the requested user or a list of all users.
//   - 404 Not Found: If a specific user is requested but not found.
//   - 500 Internal Server Error: If there is an issue fetching users from the database.
func (uc *UsersController) GetUsers(c *gin.Context) {
	/* Retrieve the user_id query parameter */
	userIDStr := c.Query("user_id")

	/* If user_id is provided, fetch the specific user */
	if userIDStr != "" {
		/* Convert user_id to uuid.UUID */
		parsedUserID, ok := utils.ParseUserIDFromString(c, userIDStr)

		if !ok {
			return /* Error already handled in the helper */
		}

		/* Query the database for the specific user */
		var user models.User
		if err := uc.DB.First(&user, "id = ?", parsedUserID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		/* Return the specific user */
		c.JSON(http.StatusOK, gin.H{"user": user})
		return
	}

	/* If no user_id is provided, fetch all users */
	var users []models.User
	if err := uc.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}

	/* Return all users */
	c.JSON(http.StatusOK, gin.H{"users": users})
}
