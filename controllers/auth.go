package controllers

import (
	"fmt"
	"net/http"
	"time"

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

// ForgotPassword handles the process of initiating a password reset for a user.
// It expects a JSON payload containing the user's email address.
//
// @param c *gin.Context - The Gin context for the HTTP request.
//
// The function performs the following steps:
// 1. Validates the input JSON to ensure the email field is present and properly formatted.
// 2. Searches for a user in the database with the provided email address.
//    - If the email does not exist, it responds with a generic success message to avoid
//      revealing whether the email is registered (security best practice).
// 3. Generates a secure password reset token and sets an expiration time for the token.
// 4. Saves the token and expiration time to the user's record in the database.
// 5. Sends a password reset email to the user containing the reset token.
//    - If email sending fails, it responds with an internal server error.
// 6. Responds with a success message indicating that a reset link has been sent.
//
// Note: The actual implementation of token generation and email sending is handled
// by utility functions (e.g., utils.GenerateRandomToken and utils.SendPasswordResetEmail).
func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	/* Find user by email */
	var user models.User
	if err := ac.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		/* Don't reveal if email doesn't exist (security best practice) */
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
		return
	}

	/* Generate reset token (use JWT or crypto/rand) */
	token := utils.GenerateRandomToken(32)
	expiresAt := time.Now().Add(1 * time.Hour)

	/* Save token to DB */
	user.PasswordResetToken = token
	user.PasswordResetExpires = expiresAt
	ac.DB.Save(&user)

	/* Send email */
	if err := utils.SendPasswordResetEmail(user.Email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset link sent to email"})
}

// ResetPassword handles the password reset process for a user.
// It expects a JSON payload with the following fields:
// - token: A valid password reset token (required).
// - new_password: The new password to set for the user (required).
//
// The function performs the following steps:
// 1. Validates the input JSON payload.
// 2. Finds the user associated with the provided token, ensuring the token is valid and not expired.
// 3. Updates the user's password after hashing it securely.
// 4. Clears the password reset token and its expiration date from the user's record.
// 5. Responds with a success message if the operation is successful, or an error message otherwise.
//
// Possible HTTP responses:
// - 400 Bad Request: If the input is invalid or the token is invalid/expired.
// - 500 Internal Server Error: If hashing the password fails.
// - 200 OK: If the password reset is successful.
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var input struct {
		Token	string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* Find user by valid token */
	var user models.User
	err := ac.DB.Where(
		"password_reset_token = ? AND password_reset_expires > ?",
		input.Token,
		time.Now(),
	).First(&user).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid/expired token"})
		return
	}

	/* Update Password */
	user.Password = input.NewPassword
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	/* Clear reset token */
	user.PasswordResetToken = ""
	user.PasswordResetExpires = time.Time{}
	ac.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

// Logout handles the user logout process.
// Since JWT is stateless, this function assumes that client-side token deletion
// is sufficient for logging out. If additional security is required, such as
// invalidating tokens server-side, a token blacklist mechanism can be implemented.
// Responds with a success message upon logout.
func (ac *AuthController) Logout(c *gin.Context) {
	/* JWT is stateless, so client-side token deletion is sufficient. */
	/* TODO:Implement a token blacklist if needed */
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// CheckAuth handles the authentication check for a user.
// 
// This method is invoked after the JWTAuth middleware has validated the user's token.
// It retrieves the user ID from the context, validates it, and fetches the user's details
// from the database. If the user is authenticated and exists in the database, their
// basic information is returned in the response.
//
// Parameters:
//   - c (*gin.Context): The Gin context containing the request and response objects.
//
// Behavior:
//   - If the user ID is not found in the context, it responds with HTTP 401 Unauthorized.
//   - If the user ID cannot be parsed into a valid UUID, it terminates the request.
//   - If the user is not found in the database, it responds with HTTP 404 Not Found.
//   - If the user is authenticated and exists, it responds with HTTP 200 OK and includes
//     the user's email, first name, and last name in the response.
//
// Response:
//   - HTTP 200 OK: Returns a JSON object with `isAuthenticated` set to true and user details.
//   - HTTP 401 Unauthorized: Returns a JSON object with an error message "Unauthorized".
//   - HTTP 404 Not Found: Returns a JSON object with an error message "User not found".
func (ac *AuthController) CheckAuth(c *gin.Context) {
	/* If JWTAuth middleware passed, the user is authenticated */
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	/* Convert userID to uuid */
	parsedUserID, ok := utils.ParseUserIDFromString(c, userID.(string))

	if !ok {
		return
	}


	var user models.User
	if err := ac.DB.Select("id, email, first_name, last_name").First(&user, "id = ?", parsedUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isAuthenticated": true,
		"user": gin.H{
			"email": user.Email,
			"firstName": user.FirstName,
			"lastName": user.LastName,
		},
	})
}