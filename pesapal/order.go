package pesapal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderRequest struct {
	ID             string  `json:"id"`
	Currency       string  `json:"currency"`
	Amount         float64 `json:"amount"`
	Description    string  `json:"description"`
	CallbackURL    string  `json:"callback_url"`
	NotificationID string  `json:"notification_id"`
	BillingAddress Address `json:"billing_address"`
}

type Address struct {
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
	Line1        string `json:"line_1"`
	Line2        string `json:"line_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	ZipCode      string `json:"zip_code"`
	CountryCode  string `json:"country_code"`
}

type OrderResponse struct {
	OrderTrackingID string `json:"order_tracking_id"`
	RedirectURL     string `json:"redirect_url"`
}

func (c *Config) SubmitOrder(authToken, ipnID string, order OrderRequest) (*OrderResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/Transactions/SubmitOrderRequest", c.BaseURL)

	/* Set the notification ID from IPN registration */
	order.NotificationID = ipnID
	order.CallbackURL = c.CallbackURL

	jsonBody, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("order submission failed with status: %s", resp.Status)
	}

	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &orderResp, nil
}
