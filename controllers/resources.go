package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type ResourceController struct {
	DB *gorm.DB
}

/*
Global cache instance
Specify the default expiration time for cached items and
the interval for purging expired items
*/
var resourceCache = cache.New(720*time.Minute, 720*time.Minute)

/* Response struct without ExtractedContent */
type ResourceResponse struct {
	ID                      uuid.UUID `json:"id"`
	ParentURL               string    `json:"parent_url"`
	GoogleDriveDownloadLink string    `json:"google_drive_download_link"`
	Name                    string    `json:"name"`
	RelativePath            string    `json:"relative_path"`
	ParentDirectory         string    `json:"parent_directory"`
	DjangoRelativePath      string    `json:"django_relative_path"`
	GoogleCloudStorageLink  string    `json:"google_cloud_storage_link"`
	CreatedAt               time.Time `json:"created_at"`
	// Categories []string `json:"categories"`
	// IsExtracted bool `json:"is_extracted"`
}

func (rc *ResourceController) GetResources(c *gin.Context) {
	var resources []models.WebCrawlerResource
	var response []ResourceResponse

	/* Generate a cache key based on search params and pagination */
	searchParams := []string{"q1", "q2", "q3", "q4"}
	cacheKey := "resources:"
	for _, param := range searchParams {
		cacheKey += param + "=" + c.Query(param) + "&"
	}
	cacheKey += "page=" + c.DefaultQuery("page", "1") + "&"
	cacheKey += "limit=" + c.DefaultQuery("limit", "100")

	/* Check if the result is already in the cache */
	if cachedData, found := resourceCache.Get(cacheKey);
	found {
		/* Return the cached response */
		c.JSON(http.StatusOK, cachedData)
		return
	}

	/* Build the query */
	query := rc.DB.Model(&models.WebCrawlerResource{})

	// searchParams := []string{"q1", "q2", "q3", "q4"}
	for _, param := range searchParams {
		if value := c.Query(param); value != "" {
			value = strings.ToLower(value)
			query = query.Where(
				"LOWER(name) LIKE ? OR "+
					"LOWER(parent_url) LIKE ? OR "+
					"LOWER(google_drive_download_link) LIKE ? OR "+
					"LOWER(relative_path) LIKE ? OR "+
					"LOWER(extracted_content) LIKE ? OR "+
					"LOWER(google_cloud_storage_link) LIKE ?",
				"%"+value+"%",
				"%"+value+"%",
				"%"+value+"%",
				"%"+value+"%",
				"%"+value+"%",
				"%"+value+"%",
			)
		}
	}

	/* Parse pagination parameters */
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "100")

	/* Count total records */
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count resources"})
		return
	}

	/* Apply pagination */
	query = query.Scopes(Paginate(page, limit))

	/* Execute query */
	if err := query.Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}

	/* Calculate pagination metadata */
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	totalPages := int((totalRecords + int64(limitInt) - 1) / int64(limitInt)) // Ceiling division
	nextPage := pageInt + 1
	if nextPage > totalPages {
		nextPage = 0 // No next page
	}

	/* Convert to response format (excluding Extracted Content) */
	// var response []ResourceResponse
	for _, r := range resources {
		response = append(response, ResourceResponse{
			ID:                      r.ID,
			ParentURL:               r.ParentURL,
			GoogleDriveDownloadLink: r.GoogleDriveDownloadLink,
			Name:                    r.Name,
			RelativePath:            r.RelativePath,
			ParentDirectory:         r.ParentDirectory,
			DjangoRelativePath:      r.DjangoRelativePath,
			GoogleCloudStorageLink:  r.GoogleCloudStorageLink,
			CreatedAt:               r.CreatedAt,
			// Categories:              r.Categories,
			// IsExtracted:             r.IsExtracted,
		})
	}

	/* Prepare the final response */
	finalResponse := gin.H{
		"data": response,
		"pagination": gin.H{
			"total_records": totalRecords,
			"total_pages": int((totalRecords + int64(limitInt) - 1) / int64(limitInt)),
			"current_page": pageInt,
			"next_page": nextPage,
			"limit": limitInt,
		},
	}

	/* Store the result in the cache */
	resourceCache.Set(cacheKey, finalResponse, cache.DefaultExpiration)

	/* Return the response */
	c.JSON(http.StatusOK, finalResponse)

	// c.JSON(http.StatusOK, gin.H{"data": response})
	/* Return response with pagination metadata */
	// c.JSON(http.StatusOK, gin.H{
	// 	"data": response,
	// 	"pagination": gin.H{
	// 		"total_records": totalRecords,
	// 		"total_pages":   totalPages,
	// 		"current_page":  pageInt,
	// 		"next_page":     nextPage,
	// 		"limit":         limitInt,
	// 	},
	// })
}

/* Pagination scope */
func Paginate(page, limit string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageInt := 1
		limitInt := 10

		// Parse page and limit from strings to integers
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageInt = p
		}
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			limitInt = l
		}

		offset := (pageInt - 1) * limitInt
		return db.Offset(offset).Limit(limitInt)
	}
}

func (rc *ResourceController) GetUniqeParentDirectories(c *gin.Context) {
	var directories []string

	/* Parse pagination parameters */
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "100")

	// Apply search query if provided
	query := rc.DB.Model(&models.WebCrawlerResource{}).Distinct("parent_directory")
	if search := c.Query("search"); search != "" {
		search = strings.ToLower(search)
		query = query.Where("LOWER(parent_directory) LIKE ?", "%"+search+"%")
	}

	// Count total records
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count directories"})
		return
	}

	/* Apply pagination */
	// query := rc.DB.Model(&models.WebCrawlerResource{}).Distinct("parent_directory")
	query = query.Scopes(Paginate(page, limit))

	/* Fetch paginated data */
	if err := query.Pluck("parent_directory", &directories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch directories"})
		return
	}

	/* Extract last segment after the last / */
	uniqueSegments := make(map[string]bool)
	var result []string

	for _, dir := range directories {
		/* Handle empty string */
		if dir == "" {
			continue
		}

		/* Remove trailing slash if exists */
		cleanedDir := strings.TrimRight(dir, "/")

		/* Split by / and get last segment */
		parts := strings.Split(cleanedDir, "/")
		lastSegment := parts[len(parts)-1]

		/* Add to result if not already present */
		if _, exists := uniqueSegments[lastSegment]; !exists {
			uniqueSegments[lastSegment] = true
			result = append(result, lastSegment)
		}
	}

	/* Calculate pagination metadata */
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	totalPages := int((totalRecords + int64(limitInt) - 1) / int64(limitInt)) // Ceiling division
	nextPage := pageInt + 1
	if nextPage > totalPages {
		nextPage = 0 // No next page
	}

	/* Return response with pagination metadata */
	c.JSON(http.StatusOK, gin.H{
		"data": result,
		"pagination": gin.H{
			"total_records": totalRecords,
			"total_pages":   totalPages,
			"current_page":  pageInt,
			"next_page":     nextPage,
			"limit":         limitInt,
		},
	})
}
