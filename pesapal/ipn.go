package pesapal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type IPNRequest struct {
	URL   string `json:"url"`
	IPNID string `json:"ipn_notification_type"`
}

type IPNResponse struct {
	URL string `json:"url"`
	ID  string `json:"ipn_id"`
}

func (c *Config) RegisterIPN(authToken string) (*IPNResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/URLSetup/RegisterIPN", c.BaseURL)

	requestBody := IPNRequest{
		URL:   c.IPNURL,
		IPNID: "GET",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
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
		return nil, fmt.Errorf("IPN registration failed with status: %s", resp.Status)
	}

	var ipnResp IPNResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipnResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &ipnResp, nil
}
