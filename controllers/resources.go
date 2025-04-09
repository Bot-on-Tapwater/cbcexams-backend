package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
)

type WebCrawlerResourceController struct {
	DB *gorm.DB
}

// NewWebCrawlerResourceController creates a new instance of the controller
func NewWebCrawlerResourceController(db *gorm.DB) *WebCrawlerResourceController {
	return &WebCrawlerResourceController{DB: db}
}

func (ctrl *WebCrawlerResourceController) GetAllWebCrawlerResources(c *gin.Context) {
	var resources []models.WebCrawlerResource

	// Get pagination parameters from query
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Convert page and limit to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	// Calculate offset
	offset := (pageInt - 1) * limitInt

	// Get total count of records
	var totalCount int64
	if err := ctrl.DB.Model(&models.WebCrawlerResource{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total count"})
		return
	}

	// Fetch paginated records
	if err := ctrl.DB.Limit(limitInt).Offset(offset).Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}

	// Calculate total pages
	totalPages := (totalCount + int64(limitInt) - 1) / int64(limitInt) // Ceiling division

	// Return the records as JSON
	c.JSON(http.StatusOK, gin.H{
		"data":        resources,
		"page":        pageInt,
		"limit":       limitInt,
		"total_pages": totalPages,
		"total_count": totalCount,
	})
}
