package httpClient

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPClient struct {
	client *http.Client
	config *HTTPClientConfig
}

func NewHTTPClient(httpClientConfig *HTTPClientConfig) *HTTPClient {
	httpClient := &http.Client{
		Timeout: httpClientConfig.Timeout * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    httpClientConfig.MaxIdleConns,
			IdleConnTimeout: httpClientConfig.IdleConnTimeout * time.Second,
		},
	}

	return &HTTPClient{client: httpClient, config: httpClientConfig}
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	for i := 0; i < c.config.MaxRetries; i++ {
		resp, err := c.client.Get(url)
		if err != nil {
			fmt.Printf("Error in connection. Try:  %d: %v\n", i+1, err)
			time.Sleep(c.config.RetryWaitDuration * time.Second)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Could not connect after %d retries", c.config.MaxRetries)
}
