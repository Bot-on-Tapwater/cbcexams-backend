package pesapal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TransactionStatusResponse struct {
	PaymentMethod            string     `json:"payment_method"`
	Amount                   float64    `json:"amount"`
	CreatedDate              CustomTime `json:"created_date"`
	ConfirmationCode         string     `json:"confirmation_code"`
	PaymentStatusDescription string     `json:"payment_status_description"`
	Description              string     `json:"description"`
	Message                  string     `json:"message"`
}

type CustomTime struct {
	time.Time
}

/* UnmarshalJSON parses the custom time format without a timezone offset. */
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Remove quotes from the JSON string
	s := string(b)
	s = s[1 : len(s)-1]

	// Parse the time without a timezone offset
	parsedTime, err := time.Parse("2006-01-02T15:04:05.000", s)
	if err != nil {
		return err
	}

	ct.Time = parsedTime
	return nil
}

func (c *Config) GetTransactionStatus(authToken, orderTrackingID string) (*TransactionStatusResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/Transactions/GetTransactionStatus?orderTrackingId=%s", c.BaseURL, orderTrackingID)

	req, err := http.NewRequest("GET", url, nil)
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
		return nil, fmt.Errorf("transaction status check failed with status: %s", resp.Status)
	}

	var statusResp TransactionStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &statusResp, nil
}
