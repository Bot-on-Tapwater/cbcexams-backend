package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/bot-on-tapwater/cbcexams-backend/pesapal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentController struct {
	PesaPalConfig *pesapal.Config
}

func NewPaymentController() *PaymentController {
	return &PaymentController{
		PesaPalConfig: pesapal.NewConfig(),
	}
}

func (pc *PaymentController) InitiatePayment(c *gin.Context) {
	/* Authenticate with PesaPal */
	authToken, err := pc.PesaPalConfig.Authenticate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed", "details": err.Error()})
		return
	}

	/* Register IPN URL (only needs to be done once, could be moved to app startup) */
	ipnResp, err := pc.PesaPalConfig.RegisterIPN(authToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "IPN registration failed", "details": err.Error()})
		return
	}

	// Set default values for the OrderRequest
	orderReq := pesapal.OrderRequest{
		ID:          uuid.NewString(),      // Generate a new UUID for the ID
		Currency:    "KES",                 // Default currency
		Description: "Default description", // Default description
		BillingAddress: pesapal.Address{
			EmailAddress: "default@example.com", // Default email
			PhoneNumber:  "+1234567890",         // Default phone number
			FirstName:    "Default",             // Default first name
			MiddleName:   "Middle",              // Default middle name
			LastName:     "User",                // Default last name
			Line1:        "123 Default Street",  // Default address line 1
			Line2:        "Suite 456",           // Default address line 2
			City:         "Default City",        // Default city
			State:        "Default State",       // Default state
			PostalCode:   "12345",               // Default postal code
			ZipCode:      "12345",               // Default zip code
			CountryCode:  "US",                  // Default country code
		},
	}

	/* Parse user-provided data from the request body */
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	/* Submit order to PesaPal */
	orderResp, err := pc.PesaPalConfig.SubmitOrder(authToken, ipnResp.ID, orderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order submission failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment initiated successfully",
		"data": gin.H{
			"order_tracking_id": orderResp.OrderTrackingID,
			"redirect_url":      orderResp.RedirectURL,
		},
	})
}

func (pc *PaymentController) CheckPaymentStatus(c *gin.Context) {
	orderTrackingID := c.Param("order_id")
	if orderTrackingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order tracking ID is required"})
		return
	}

	authToken, err := pc.PesaPalConfig.Authenticate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed", "details": err.Error()})
		return
	}

	statusResp, err := pc.PesaPalConfig.GetTransactionStatus(authToken, orderTrackingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment status retrieved",
		"data":    statusResp,
	})
}

func (pc *PaymentController) HandleIPN(c *gin.Context) {
	/* Extract query parameters sent by Pesapal */
	ipnData := c.Request.URL.Query()

	/* Convert query parameters to a map for easier handling */
	data := make(map[string]string)
	for key, values := range ipnData {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}

	/* Define the JSON file path */
	filePath := "/home/bot-on-tapwater/cbcexams-backend/ipn_notifications.json"

	/* Read existing data from the file (if it exists) */
	var existingData []map[string]string
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open IPN file", "details": err.Error()})
		return
	}
	defer file.Close()

	/* Decode existing JSON data (if any) */
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&existingData); err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read IPN file", "details": err.Error()})
		return
	}

	/* Append the new data to the existing data */
	existingData = append(existingData, data)

	/* Write the updated data back to the file */
	file.Truncate(0) /* Clear the file */
	file.Seek(0, 0)  /* Move the pointer to the beginning */
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(existingData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write IPN file", "details": err.Error(),
		})
		return
	}

	/* Respond to PesaPal */
	c.JSON(http.StatusOK, gin.H{"message": "IPN received and stored successfully"})
}
