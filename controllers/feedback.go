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
		"data": input,
	})	
}

/* Get all feedback (for admin) */
func (fc *FeedbackController) GetFeedback(c *gin.Context) {
	var feedback []models.Feedback

	if err := fc.DB.Find(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedback})
}