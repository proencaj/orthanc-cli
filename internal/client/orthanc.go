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

	// Get the current context configuration
	orthancCfg, err := cfg.GetCurrentContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get current context: %w", err)
	}

	// Validate required configuration
	if orthancCfg.URL == "" {
		return nil, fmt.Errorf("orthanc URL is required (use 'orthanc config set-context %s --url <url>')", cfg.CurrentContext)
	}

	// Create client options
	var opts []gorthanc.ClientOption

	// Add basic authentication if credentials are provided
	if orthancCfg.Username != "" && orthancCfg.Password != "" {
		opts = append(opts, gorthanc.WithBasicAuth(orthancCfg.Username, orthancCfg.Password))
	}

	// Create custom HTTP client if insecure mode is enabled
	if orthancCfg.Insecure {
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
	client, err := gorthanc.NewClient(orthancCfg.URL, opts...)
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
	orthancCfg, err := c.config.GetCurrentContext()
	if err != nil {
		return ""
	}
	return orthancCfg.URL
}
