package controllers

import (
	"fmt"
	"net/http"

	"github.com/bot-on-tapwater/cbcexams-backend/models"
	"github.com/bot-on-tapwater/cbcexams-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

/* Register a new user */
// Register handles the user registration process.
// It binds the incoming JSON payload to a User model, hashes the user's password,
// and saves the user to the database. If any step fails, it returns an appropriate
// HTTP error response.
//
// @param c *gin.Context - The Gin context containing the request and response objects.
//
// Possible Responses:
// - 400 Bad Request: If the JSON payload is invalid.
// - 409 Conflict: If the email already exists in the database.
// - 500 Internal Server Error: If password hashing fails.
// - 201 Created: If the user is successfully created.
func (ac *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Password before hashing: %s\n", user.Password)
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

/* Login and return JWT */
// Login handles user authentication by validating the provided email and password.
// It expects a JSON payload with "email" and "password" fields, both of which are required.
// If the credentials are valid, it generates a JWT token and returns it in the response.
//
// @Summary      User Login
// @Description  Authenticates a user and returns a JWT token upon successful login.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body  struct{Email string; Password string}  true  "User credentials"
// @Success      200  {object}  map[string]string  "JWT token"
// @Failure      400  {object}  map[string]string  "Bad request error"
// @Failure      404  {object}  map[string]string  "User not found error"
// @Failure      401  {object}  map[string]string  "Invalid credentials error"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var credentials struct {
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := ac.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := user.CheckPassword(credentials.Password); err != nil {
		fmt.Println(credentials.Password)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"token":token})
}