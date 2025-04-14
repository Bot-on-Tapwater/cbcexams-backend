package pesapal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthResponse struct {
	Token      string `json:"token"`
	ExpiryDate string `json:"expiryDate"`
}

func (c *Config) Authenticate() (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/Auth/RequestToken", c.BaseURL)

	/* Create the request body with consumer_key and consumer_secret */
	requestBody := map[string]string{
		"consumer_key":    c.ConsumerKey,
		"consumer_secret": c.ConsumerSecret,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	// req.Header.Add("Authorization", c.GetAuthHeader())

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return "", fmt.Errorf("authentication failed with status: %s", resp.Status)
	// }

	if resp.StatusCode != http.StatusOK {
		// Read the response body for more details
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("authentication failed with status: %s, response: %s", resp.Status, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	return authResp.Token, nil
}
