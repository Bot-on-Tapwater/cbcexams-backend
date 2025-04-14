package controllers

import (
	"net/http"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookmarkController struct {
	DB *gorm.DB
}

/* Bookmark a resource */
func (bc *BookmarkController) CreateBookmark(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id")) /* From JWT middleware */
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var input struct {
		ResourceID string `json:"resource_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceID, err := uuid.Parse(input.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	/* Check if resource exists */
	var resource models.WebCrawlerResource
	if err := bc.DB.First(&resource, "id = ?", resourceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	/* Check for existing bookmark */
	var existing models.Bookmark
	if err := bc.DB.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Resource already bookmarked"})
		return
	}

	bookmark := models.Bookmark{
		UserID:     userID,
		ResourceID: resourceID,
	}

	if err := bc.DB.Create(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Resource bookmarked successfully",
		"data": gin.H{
			"bookmark_id": bookmark.ID,
			"resource":    resource,
		},
	})
}

/* Remove bookmark */
func (bc *BookmarkController) DeleteBookmark(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	resourceID, err := uuid.Parse(c.Param("resource_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	result := bc.DB.Where("user_id = ? AND resource_id = ?", userID, resourceID).Delete(&models.Bookmark{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove bookmark"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark removed successfully"})
}

/* Get all bookmarks for current user */
func (bc *BookmarkController) GetUserBookmarks(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	var bookmarks []models.Bookmark
	err = bc.DB.Preload("Resource").Where("user_id = ?", userID).Order("created_at DESC").Find(&bookmarks).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookmarks"})
		return
	}

	/* Extract resources from bookmarks */
	var resources []models.WebCrawlerResource
	for _, b := range bookmarks {
		resources = append(resources, b.Resource)
	}

	c.JSON(http.StatusOK, gin.H{"data": resources})
}
