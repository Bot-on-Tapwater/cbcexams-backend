package controllers

import (
	"net/http"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type JobsController struct {
	DB *gorm.DB
}

/* School Job Listings */
// CreateSchoolJob handles the creation of a new school job listing.
// @Summary Create a new school job listing
// @Description This endpoint allows clients to create a new school job listing by providing the necessary details in JSON format.
// @Tags Jobs
// @Accept json
// @Produce json
// @Param input body models.SchoolJobListing true "School Job Listing Input"
// @Success 201 {object} gin.H{"data": models.SchoolJobListing}
// @Failure 400 {object} gin.H{"error": string} "Bad Request"
// @Failure 500 {object} gin.H{"error": string} "Internal Server Error"
// @Router /jobs/school [post]
func (jc *JobsController) CreateSchoolJob(c *gin.Context) {
	var input models.SchoolJobListing
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := jc.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job listing"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// GetSchoolJobs handles the HTTP request to retrieve a list of school job listings.
// It supports optional query parameters for filtering the results by subject and location.
//
// Query Parameters:
//   - subject: (optional) Filters the job listings by the specified subject.
//   - location: (optional) Filters the job listings by the specified location.
//
// Response:
//   - 200 OK: Returns a JSON object containing the list of job listings.
//     Example: {"data": [{"id": 1, "subject": "Math", "location": "New York", ...}, ...]}
//   - 500 Internal Server Error: Returns a JSON object with an error message if the query fails.
//     Example: {"error": "Failed to fetch jobs"}
//
// This function expects the database connection to be available in the JobsController (jc.DB)
// and uses the Gin framework for handling HTTP requests and responses.
func (jc *JobsController) GetSchoolJobs(c *gin.Context) {
	var jobs []models.SchoolJobListing
	query := jc.DB.Where("is_active = ?", true)

	/* Optional filters */
	if subject := c.Query("subject"); subject != "" {
		query = query.Where("subject = ?", subject)
	}
	if location := c.Query("location"); location != "" {
		query = query.Where("location = ?", location)
	}

	if err := query.Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": jobs})
}

/* Teacher Profiles */
// CreateTeacherProfile handles the creation of a teacher's job profile.
// It binds form data, processes multi-value fields, handles file uploads,
// and saves the profile to the database.
//
// @param c *gin.Context - The Gin context containing the HTTP request and response.
//
// The function performs the following steps:
// 1. Binds form data to a TeacherJobProfile model.
// 2. Parses multi-value form fields for "subjects" and "education_level".
// 3. Handles the upload of a resume file and saves it to the server.
// 4. Saves the teacher profile to the database.
//
// Responses:
// - 400 Bad Request: If the form data binding fails.
// - 500 Internal Server Error: If file upload or database operations fail.
// - 201 Created: If the profile is successfully created.
//
// Expected form fields:
// - "subjects" (multi-value): The subjects the teacher specializes in.
// - "education_level" (multi-value): The education levels the teacher can teach.
// - "resume" (file): The teacher's resume file.
func (jc *JobsController) CreateTeacherProfile(c *gin.Context) {
	var input models.TeacherJobProfile

	/* Bind form data (supports multipart/form-data) */
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* Parse multi-value form fields manually */
	input.Subjects = pq.StringArray(c.PostFormArray("subjects"))
	input.EducationLevel = pq.StringArray(c.PostFormArray("education_level"))

	/* Handle resume upload */
	if file, err := c.FormFile("resume"); err == nil {
		resumePath := "uploads/resumes/" + file.Filename
		if err := c.SaveUploadedFile(file, resumePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume"})
			return
		}
		input.ResumePath = resumePath
	}

	if err := jc.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// GetTeacherProfiles retrieves a list of active teacher profiles from the database
// based on optional query parameters for filtering by subject and location.
//
// Query Parameters:
//   - subject: (optional) Filters teachers by a specific subject. Matches if the subject
//     is present in the teacher's list of subjects.
//   - location: (optional) Filters teachers by their location.
//
// Response:
//   - On success: Returns a JSON object with an HTTP status of 200 containing the list
//     of teacher profiles in the "data" field.
//   - On failure: Returns a JSON object with an HTTP status of 500 containing an "error"
//     field with a failure message.
func (jc *JobsController) GetTeacherProfiles(c *gin.Context) {
	var teachers []models.TeacherJobProfile
	query := jc.DB.Where("is_active = ?", true)

	/* Optional filters */
	if subject := c.Query("subject"); subject != "" {
		query = query.Where("? = ANY(subjects)", subject)
	}
	if location := c.Query("location"); location != "" {
		query = query.Where("location = ?", location)
	}

	if err := query.Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teachers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": teachers})
}
