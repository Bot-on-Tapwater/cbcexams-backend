package controllers

import (
	"net/http"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TutoringController struct {
	DB *gorm.DB
}

/* Tutor Requests (Students/Parents) */
// CreateTutorRequest handles the creation of a new tutor request.
// It expects a JSON payload in the request body that matches the TutorRequest model.
// If the payload is invalid, it responds with a 400 Bad Request status and an error message.
// If the database operation fails, it responds with a 500 Internal Server Error status and an error message.
// On success, it responds with a 201 Created status and the created tutor request data.
//
// @param c *gin.Context - The Gin context containing the HTTP request and response.
// @response 400 - If the JSON payload is invalid.
// @response 500 - If there is an error while saving the request to the database.
// @response 201 - If the tutor request is successfully created.
func (tc *TutoringController) CreateTutorRequest(c *gin.Context) {
	var input models.TutorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// GetTutorRequests handles the HTTP GET request to retrieve all tutor requests.
// It queries the database for all records in the TutorRequest model and returns
// them as a JSON response. If an error occurs during the database query, it
// responds with an HTTP 500 status and an error message.
//
// @Summary Retrieve tutor requests
// @Description Fetches all tutor requests from the database
// @Tags Tutoring
// @Produce json
// @Success 200 {object} gin.H{"data": []models.TutorRequest}
// @Failure 500 {object} gin.H{"error": string}
// @Router /tutor-requests [get]
func (tc *TutoringController) GetTutorRequests(c *gin.Context) {
	var requests []models.TutorRequest
	if err := tc.DB.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requests})
}

/* Tutor applications (Tutors) */
// CreateTutorApplication handles the creation of a new tutor application.
// It expects JSON data for the main application fields and supports additional
// multi-value form fields and file uploads.
//
// @param c *gin.Context - The Gin context containing the request and response objects.
//
// Request Body:
// - JSON fields for the TutorApplication model.
//
// Form Fields:
// - "subjects" (array): A list of subjects the tutor is applying to teach.
// - "education_level" (array): The tutor's education levels.
// - "available_days" (array): Days the tutor is available.
//
// File Upload:
// - "resume" (optional): A file upload for the tutor's resume.
//
// Responses:
// - 201 Created: Returns the created TutorApplication object.
// - 400 Bad Request: If the input data is invalid.
// - 500 Internal Server Error: If there is an error saving the file or creating the application.
func (tc *TutoringController) CreateTutorApplication(c *gin.Context) {
	var input models.TutorApplication

	/* Bind JSON data */
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse multi-value form fields manually
	input.Subjects = pq.StringArray(c.PostFormArray("subjects"))
	input.EducationLevel = pq.StringArray(c.PostFormArray("education_level"))
	input.AvailableDays = pq.StringArray(c.PostFormArray("available_days"))

	/* Handle file upload (optional) */
	if file, err := c.FormFile("resume"); err == nil {
		resumePath := "uploads/resumes/" + file.Filename
		if err := c.SaveUploadedFile(file, resumePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume"})
			return
		}
		input.ResumePath = resumePath
	}

	if err := tc.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// GetTutorApplications retrieves all tutor applications from the database
// and returns them as a JSON response. If an error occurs during the database
// query, it responds with an HTTP 500 status and an error message.
//
// @Summary Retrieve tutor applications
// @Description Fetches all tutor applications from the database.
// @Tags Tutoring
// @Produce json
// @Success 200 {object} gin.H{"data": []models.TutorApplication}
// @Failure 500 {object} gin.H{"error": string}
// @Router /tutor-applications [get]
func (tc *TutoringController) GetTutorApplications(c *gin.Context) {
	var applications []models.TutorApplication
	if err := tc.DB.Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": applications})
}
