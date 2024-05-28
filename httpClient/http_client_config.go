package httpClient

import (
	"fmt"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

type HTTPClientConfig struct {
	Timeout           time.Duration `yaml:"timeout"`
	MaxIdleConns      int           `yaml:"max_idle_conns"`
	IdleConnTimeout   time.Duration `yaml:"idle_conn_timeout"`
	MaxRetries        int           `yaml:"max_retries"`
	RetryWaitDuration time.Duration `yaml:"retry_wait_duration"`
}

func NewHTTPClientConfig(configFile string) (*HTTPClientConfig, error) {
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	var httpClientConfig *HTTPClientConfig
	if err := yaml.Unmarshal(configData, &httpClientConfig); err != nil {
		return nil, fmt.Errorf("Error decoding config file: %v", err)
	}
	println("This is the config")
	println(httpClientConfig.MaxRetries)
	return httpClientConfig, nil
}
