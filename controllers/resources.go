package controllers

import (
	"net/http"
	"regexp"
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
	ID uuid.UUID `json:"id"`
	// ParentURL               string    `json:"parent_url"`
	// GoogleDriveDownloadLink string    `json:"google_drive_download_link"`
	Name string `json:"name"`
	// RelativePath            string    `json:"relative_path"`
	// ParentDirectory         string    `json:"parent_directory"`
	DjangoRelativePath     string    `json:"django_relative_path"`
	GoogleCloudStorageLink string    `json:"google_cloud_storage_link"`
	CreatedAt              time.Time `json:"created_at"`
	// Categories []string `json:"categories"`
	// IsExtracted bool `json:"is_extracted"`
}

// Helper function to add space after "form" or "grade" if followed by a number without space
func addSpaceAfterFormOrGrade(input string) string {
	// Check for "form" followed by a number without space
	reForm := regexp.MustCompile(`(?i)form(\d+)`)
	if reForm.MatchString(input) {
		input = reForm.ReplaceAllString(input, "form $1")
	}

	// Check for "grade" followed by a number without space
	reGrade := regexp.MustCompile(`(?i)grade(\d+)`)
	if reGrade.MatchString(input) {
		input = reGrade.ReplaceAllString(input, "grade $1")
	}

	return input
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
	if cachedData, found := resourceCache.Get(cacheKey); found {
		/* Return the cached response */
		c.JSON(http.StatusOK, cachedData)
		return
	}

	/* Parse pagination parameters */
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	/* Try different parameter combinations until we find results */
	var finalQuery *gorm.DB
	var totalRecords int64
	var queryUsed []string // Track which parameters were used

	// Try all parameters first, then fall back to fewer parameters if no results
	for i := len(searchParams); i > 0; i-- {
		query := rc.DB.Model(&models.WebCrawlerResource{})
		queryUsed = searchParams[:i] // Current parameters being tried
		hasConditions := false

		// Apply all parameters up to the current index
		for _, param := range queryUsed {
			if value := c.Query(param); value != "" {
				value = strings.ToLower(value)
				// Apply form/grade transformation only for q1
				if param == "q1" {
					value = addSpaceAfterFormOrGrade(value)
				}
				query = query.Where(
					"LOWER(name) LIKE ? OR "+
						"LOWER(parent_directory) LIKE ? OR "+
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
				hasConditions = true
			}
		}

		// If no conditions were applied (all params empty), break and use all records
		if !hasConditions {
			queryUsed = []string{} // No parameters used
			break
		}

		// Count records for this query
		var count int64
		if err := query.Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count resources"})
			return
		}

		// If we found results, use this query
		if count > 0 {
			finalQuery = query
			totalRecords = count
			break
		}
	}

	// If all parameter combinations returned 0 results, use base query (no conditions)
	if finalQuery == nil {
		finalQuery = rc.DB.Model(&models.WebCrawlerResource{})
		if err := finalQuery.Count(&totalRecords).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count resources"})
			return
		}
	}

	/* Apply pagination */
	finalQuery = finalQuery.Scopes(Paginate(page, limit))

	/* Execute query */
	if err := finalQuery.Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}

	/* Calculate pagination metadata */
	totalPages := int((totalRecords + int64(limitInt) - 1) / int64(limitInt)) // Ceiling division
	nextPage := pageInt + 1
	if nextPage > totalPages {
		nextPage = 0 // No next page
	}

	/* Convert to response format (excluding Extracted Content) */
	for _, r := range resources {
		response = append(response, ResourceResponse{
			ID:                     r.ID,
			Name:                   r.Name,
			DjangoRelativePath:     r.DjangoRelativePath,
			GoogleCloudStorageLink: r.GoogleCloudStorageLink,
			CreatedAt:              r.CreatedAt,
		})
	}

	/* Prepare the final response */
	finalResponse := gin.H{
		"data": response,
		"pagination": gin.H{
			"total_records": totalRecords,
			"total_pages":   totalPages,
			"current_page":  pageInt,
			"next_page":     nextPage,
			"limit":         limitInt,
		},
		"parameters_used": queryUsed, // Include which parameters were actually used
	}

	/* Store the result in the cache */
	resourceCache.Set(cacheKey, finalResponse, cache.DefaultExpiration)

	/* Return the response */
	c.JSON(http.StatusOK, finalResponse)
}

// func (rc *ResourceController) GetResources(c *gin.Context) {
// 	var resources []models.WebCrawlerResource
// 	var response []ResourceResponse

// 	/* Generate a cache key based on search params and pagination */
// 	searchParams := []string{"q1", "q2", "q3", "q4"}

// 	cacheKey := "resources:"
// 	for _, param := range searchParams {
// 		cacheKey += param + "=" + c.Query(param) + "&"
// 	}
// 	cacheKey += "page=" + c.DefaultQuery("page", "1") + "&"
// 	cacheKey += "limit=" + c.DefaultQuery("limit", "100")

// 	/* Check if the result is already in the cache */
// 	if cachedData, found := resourceCache.Get(cacheKey); found {
// 		/* Return the cached response */
// 		c.JSON(http.StatusOK, cachedData)
// 		return
// 	}

// 	/* Build the query */
// 	query := rc.DB.Model(&models.WebCrawlerResource{})

// 	for _, param := range searchParams {
// 		if value := c.Query(param); value != "" {
// 			value = strings.ToLower(value)
// 			// Apply form/grade transformation only for q1
// 			if param == "q1" {
// 				value = addSpaceAfterFormOrGrade(value)
// 			}
// 			query = query.Where(
// 				"LOWER(name) LIKE ? OR "+
// 					"LOWER(parent_directory) LIKE ? OR "+
// 					"LOWER(google_drive_download_link) LIKE ? OR "+
// 					"LOWER(relative_path) LIKE ? OR "+
// 					"LOWER(extracted_content) LIKE ? OR "+
// 					"LOWER(google_cloud_storage_link) LIKE ?",
// 				"%"+value+"%",
// 				"%"+value+"%",
// 				"%"+value+"%",
// 				"%"+value+"%",
// 				"%"+value+"%",
// 				"%"+value+"%",
// 			)
// 		}
// 	}

// 	/* Parse pagination parameters */
// 	page := c.DefaultQuery("page", "1")
// 	limit := c.DefaultQuery("limit", "100")

// 	/* Count total records */
// 	var totalRecords int64
// 	if err := query.Count(&totalRecords).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count resources"})
// 		return
// 	}

// 	/* Apply pagination */
// 	query = query.Scopes(Paginate(page, limit))

// 	/* Execute query */
// 	if err := query.Find(&resources).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
// 		return
// 	}

// 	/* Calculate pagination metadata */
// 	pageInt, _ := strconv.Atoi(page)
// 	limitInt, _ := strconv.Atoi(limit)
// 	totalPages := int((totalRecords + int64(limitInt) - 1) / int64(limitInt)) // Ceiling division
// 	nextPage := pageInt + 1
// 	if nextPage > totalPages {
// 		nextPage = 0 // No next page
// 	}

// 	/* Convert to response format (excluding Extracted Content) */
// 	for _, r := range resources {
// 		response = append(response, ResourceResponse{
// 			ID: r.ID,
// 			// ParentURL:               r.ParentURL,
// 			// GoogleDriveDownloadLink: r.GoogleDriveDownloadLink,
// 			Name: r.Name,
// 			// RelativePath:            r.RelativePath,
// 			// ParentDirectory:         r.ParentDirectory,
// 			DjangoRelativePath:     r.DjangoRelativePath,
// 			GoogleCloudStorageLink: r.GoogleCloudStorageLink,
// 			CreatedAt:              r.CreatedAt,
// 			// Categories:              r.Categories,
// 			// IsExtracted:             r.IsExtracted,
// 		})
// 	}

// 	/* Prepare the final response */
// 	finalResponse := gin.H{
// 		"data": response,
// 		"pagination": gin.H{
// 			"total_records": totalRecords,
// 			"total_pages":   int((totalRecords + int64(limitInt) - 1) / int64(limitInt)),
// 			"current_page":  pageInt,
// 			"next_page":     nextPage,
// 			"limit":         limitInt,
// 		},
// 	}

// 	/* Store the result in the cache */
// 	resourceCache.Set(cacheKey, finalResponse, cache.DefaultExpiration)

// 	/* Return the response */
// 	c.JSON(http.StatusOK, finalResponse)
// }

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
	search := strings.ReplaceAll(c.Query("search"), " ", "-")
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "100")
    cacheKey := "unique_directories:search=" + search + "&page=" + page + "&limit=" + limit

    if cachedData, found := resourceCache.Get(cacheKey); found {
        c.JSON(http.StatusOK, cachedData)
        return
    }

    // Step 1: Get distinct parent directories with pagination
    var directories []string
    query := rc.DB.Model(&models.WebCrawlerResource{}).Distinct("parent_directory")
    
    if search != "" {
        search = strings.ToLower(search)
        query = query.Where("LOWER(parent_directory) LIKE ?", "%"+search+"%")
    }

    // Count total unique directories
    var totalRecords int64
    if err := query.Count(&totalRecords).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count directories"})
        return
    }

    // Apply pagination to the directory query
    pageInt, _ := strconv.Atoi(page)
    limitInt, _ := strconv.Atoi(limit)
    offset := (pageInt - 1) * limitInt
    
    if err := query.Offset(offset).Limit(limitInt).Pluck("parent_directory", &directories).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch directories"})
        return
    }

    // Step 2: Get only the required fields for these directories
    type ResourceResponse struct {
        ID                     uuid.UUID `json:"id"`
        ParentDirectory        string    `json:"parent_directory"`
        Name                   string    `json:"name"`
        GoogleCloudStorageLink string    `json:"google_cloud_storage_link"`
    }
    
    var allRecords []ResourceResponse
    if err := rc.DB.Model(&models.WebCrawlerResource{}).
        Select("id", "parent_directory", "name", "google_cloud_storage_link").
        Where("parent_directory IN ?", directories).
        Find(&allRecords).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
        return
    }

    // Process the records
    prefix := "/home/bot-on-tapwater/projects/cbcexams/media/downloaded_files/"
    directoryMap := make(map[string][]ResourceResponse)

    for _, record := range allRecords {
        trimmedDir := strings.TrimPrefix(record.ParentDirectory, prefix)
        if trimmedDir != "" {
            directoryMap[trimmedDir] = append(directoryMap[trimmedDir], record)
        }
    }

    // Calculate pagination
    totalPages := int((totalRecords + int64(limitInt) - 1) / int64(limitInt))
    nextPage := pageInt + 1
    if nextPage > totalPages {
        nextPage = 0
    }

    finalResponse := gin.H{
        "data": directoryMap,
        "pagination": gin.H{
            "total_records": totalRecords,
            "total_pages":   totalPages,
            "current_page":  pageInt,
            "next_page":     nextPage,
            "limit":         limitInt,
        },
    }

    resourceCache.Set(cacheKey, finalResponse, cache.DefaultExpiration)
    c.JSON(http.StatusOK, finalResponse)
}