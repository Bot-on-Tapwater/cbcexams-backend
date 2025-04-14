package controllers

import (
	"net/http"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebDevController struct {
	DB *gorm.DB
}

/* Create new website request */

// CreateRequest handles the creation of a new web development request.
// It expects a JSON payload in the request body that conforms to the WebDevRequest model.
//
// The function performs the following steps:
// 1. Binds the incoming JSON payload to the WebDevRequest struct.
// 2. Validates the input, ensuring that if the client type is "school", the organization name is provided.
// 3. Attempts to save the request to the database.
// 4. Returns appropriate HTTP responses based on the success or failure of the operation.
//
// Responses:
// - 400 Bad Request: If the input JSON is invalid.
// - 500 Internal Server Error: If validation fails or the database operation encounters an error.
// - 201 Created: If the request is successfully saved.
//
// Example JSON payload:
//
//	{
//	  "ClientType": "school",
//	  "OrganizationName": "Example School",
//	  "OtherFields": "Other values specific to the WebDevRequest model"
//	}
func (wdc *WebDevController) CreateRequest(c *gin.Context) {
	var input models.WebDevRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* Validate school clients */
	if input.ClientType == "school" && input.OrganizationName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit request"})
		return
	}

	if err := wdc.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Request submitted successfully",
		"data":    input,
	})
}

/* Get all requests (with optional filters) */
// GetRequests handles the HTTP GET request to fetch web development requests.
// It supports optional query parameters for filtering the results:
// - "status": Filters requests by their status (e.g., "pending", "completed").
// - "client_type": Filters requests by the type of client (e.g., "individual", "business").
//
// The function queries the database for matching WebDevRequest records and returns
// them in the response. If an error occurs during the database query, it responds
// with an HTTP 500 status and an error message.
//
// Parameters:
// - c (*gin.Context): The Gin context, which provides access to query parameters
//   and handles the HTTP response.
//
// Response:
// - HTTP 200: Returns a JSON object containing the list of matching requests.
// - HTTP 500: Returns a JSON object with an error message if the query fails.
func (wdc *WebDevController) GetRequests(c *gin.Context) {
	var requests []models.WebDevRequest
	query := wdc.DB

	/* Optional filters */
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if clientType := c.Query("client_type"); clientType != "" {
		query = query.Where("client_type = ?", clientType)
	}

	if err := query.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requests})
}
