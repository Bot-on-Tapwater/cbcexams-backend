package controllers

import (
	"net/http"

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
