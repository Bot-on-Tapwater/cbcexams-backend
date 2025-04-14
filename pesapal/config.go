package pesapal

import (
	"encoding/base64"
	"fmt"
	"os"
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	BaseURL        string
	CallbackURL    string
	IPNURL         string
}

func NewConfig() *Config {
	return &Config{
		ConsumerKey:    os.Getenv("PESAPAL_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("PESAPAL_CONSUMER_SECRET"),
		BaseURL:        os.Getenv("PESAPAL_BASE_URL"),
		CallbackURL:    os.Getenv("PESAPAL_CALLBACK_URL"),
		IPNURL:         os.Getenv("PESAPAL_IPN_URL"),
	}
}

func (c *Config) GetAuthHeader() string {
	credentials := fmt.Sprintf("%s:%s", c.ConsumerKey, c.ConsumerSecret)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}
