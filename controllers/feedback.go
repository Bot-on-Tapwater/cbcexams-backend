package controllers

import (
	"net/http"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FeedbackController struct {
	DB *gorm.DB
}

/* Submit feedback */
// SubmitFeedback handles the submission of feedback from the client.
// It expects a JSON payload that binds to the models.Feedback struct.
// If the payload is invalid, it responds with a 400 Bad Request status and an error message.
// If the feedback cannot be saved to the database, it responds with a 500 Internal Server Error status and an error message.
// On successful submission, it responds with a 201 Created status, a success message, and the submitted feedback data.
func (fc *FeedbackController) SubmitFeedback(c *gin.Context) {
	var input models.Feedback

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fc.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit feedback"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Feedback submitted successfully",
		"data":    input,
	})
}

/* Get all feedback (for admin) */
// GetFeedback handles the HTTP GET request to retrieve all feedback records.
// It queries the database for feedback entries and returns them as a JSON response.
// If an error occurs during the database query, it responds with an HTTP 500 status
// and an error message.
//
// @Summary Retrieve all feedback
// @Description Fetches all feedback records from the database
// @Tags Feedback
// @Produce json
// @Success 200 {object} gin.H{"data": []models.Feedback}
// @Failure 500 {object} gin.H{"error": string}
// @Router /feedback [get]
func (fc *FeedbackController) GetFeedback(c *gin.Context) {
	var feedback []models.Feedback

	if err := fc.DB.Find(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedback})
}
