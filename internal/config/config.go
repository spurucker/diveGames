package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"log"
	"net/http"
	"os"
	"time"
)

type HTTPClient struct {
	client *http.Client
	config *HTTPClientConfig
}

type HTTPClientConfig struct {
	Timeout           time.Duration `yaml:"timeout"`
	MaxIdleConns      int           `yaml:"max_idle_conns"`
	IdleConnTimeout   time.Duration `yaml:"idle_conn_timeout"`
	MaxRetries        int           `yaml:"max_retries"`
	RetryWaitDuration time.Duration `yaml:"retry_wait_duration"`
}

func NewHTTPClient(configFile string) (*HTTPClient, error) {
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var httpClientConfig *HTTPClientConfig
	if err := yaml.Unmarshal(configData, &httpClientConfig); err != nil {
		return nil, fmt.Errorf("error decoding config file: %v", err)
	}

	httpClient := &http.Client{
		Timeout: httpClientConfig.Timeout * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    httpClientConfig.MaxIdleConns,
			IdleConnTimeout: httpClientConfig.IdleConnTimeout * time.Second,
		},
	}

	return &HTTPClient{client: httpClient, config: httpClientConfig}, nil
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	for i := 0; i < c.config.MaxRetries; i++ {
		resp, err := c.client.Get(url)
		if err != nil {
			log.Printf("Error in connection. Try:  %d: %v\n\n", i+1, err)
			time.Sleep(c.config.RetryWaitDuration * time.Second)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("could not connect after %d retries", c.config.MaxRetries)
}
