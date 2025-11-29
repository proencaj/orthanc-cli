package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/proencaj/gorthanc"
	"github.com/proencaj/orthanc-cli/internal/config"
)

// Client wraps the gorthanc client
type Client struct {
	*gorthanc.Client
	config *config.Config
}

// NewClient creates a new Orthanc client from the configuration
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration is required")
	}

	// Validate required configuration
	if cfg.Orthanc.URL == "" {
		return nil, fmt.Errorf("orthanc URL is required (use 'orthanc config set orthanc.url <url>')")
	}

	// Create client options
	var opts []gorthanc.ClientOption

	// Add basic authentication if credentials are provided
	if cfg.Orthanc.Username != "" && cfg.Orthanc.Password != "" {
		opts = append(opts, gorthanc.WithBasicAuth(cfg.Orthanc.Username, cfg.Orthanc.Password))
	}

	// Create custom HTTP client if insecure mode is enabled
	if cfg.Orthanc.Insecure {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		httpClient := &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}
		opts = append(opts, gorthanc.WithHTTPClient(httpClient))
	}

	// Create the gorthanc client
	client, err := gorthanc.NewClient(cfg.Orthanc.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create orthanc client: %w", err)
	}

	return &Client{
		Client: client,
		config: cfg,
	}, nil
}

// GetConfig returns the configuration used by the client
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// URL returns the Orthanc server URL
func (c *Client) URL() string {
	return c.config.Orthanc.URL
}
