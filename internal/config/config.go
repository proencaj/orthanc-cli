package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Contexts       map[string]*ContextConfig `mapstructure:"contexts"`
	CurrentContext string                    `mapstructure:"current-context"`
	Output         OutputConfig              `mapstructure:"output"`

	// Legacy fields for backward compatibility (deprecated)
	Orthanc *OrthancConfig `mapstructure:"orthanc,omitempty"`
}

// ContextConfig holds configuration for a single context
type ContextConfig struct {
	Orthanc OrthancConfig `mapstructure:"orthanc"`
}

// OrthancConfig holds Orthanc server configuration
type OrthancConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Insecure bool   `mapstructure:"insecure"`
}

// OutputConfig holds output formatting configuration
type OutputConfig struct {
	JSON bool `mapstructure:"json"`
}

// GetCurrentContext returns the configuration for the current context
// with environment variable overrides applied
func (c *Config) GetCurrentContext() (*OrthancConfig, error) {
	if c.CurrentContext == "" {
		return nil, fmt.Errorf("no context selected")
	}

	ctx, exists := c.Contexts[c.CurrentContext]
	if !exists {
		return nil, fmt.Errorf("context %q not found", c.CurrentContext)
	}

	// Make a copy to avoid modifying the original
	config := ctx.Orthanc

	// Apply environment variable overrides
	// Environment variables follow the pattern: ORTHANC_URL, ORTHANC_USERNAME, etc.
	if url := viper.GetString("url"); url != "" {
		config.URL = url
	}
	if username := viper.GetString("username"); username != "" {
		config.Username = username
	}
	if password := viper.GetString("password"); password != "" {
		config.Password = password
	}
	if viper.IsSet("insecure") {
		config.Insecure = viper.GetBool("insecure")
	}

	return &config, nil
}

// LoadConfig reads configuration from file and environment variables
func LoadConfig(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		// Search config in home directory with name ".orthanc-cli" (without extension)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".orthanc-cli")
	}

	// Environment variables
	viper.SetEnvPrefix("ORTHANC")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; using defaults and environment variables
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Migrate legacy config format to multi-context format
	if err := migrateConfig(&config); err != nil {
		return nil, fmt.Errorf("error migrating config: %w", err)
	}

	return &config, nil
}

// migrateConfig migrates legacy single-context config to multi-context format
func migrateConfig(config *Config) error {
	// Check if this is a legacy config (has orthanc field but no contexts)
	if config.Orthanc != nil && len(config.Contexts) == 0 {
		// Create a default context from legacy config
		config.Contexts = make(map[string]*ContextConfig)
		config.Contexts["default"] = &ContextConfig{
			Orthanc: *config.Orthanc,
		}
		config.CurrentContext = "default"

		// Save the migrated config
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			configFile = filepath.Join(home, ".orthanc-cli.yaml")
		}

		// Clear the legacy orthanc field
		config.Orthanc = nil

		// Write migrated config
		viper.Set("contexts", config.Contexts)
		viper.Set("current-context", config.CurrentContext)
		viper.Set("output", config.Output)
		viper.Set("orthanc", nil)

		if err := viper.WriteConfigAs(configFile); err != nil {
			return fmt.Errorf("failed to write migrated config: %w", err)
		}

		fmt.Fprintf(os.Stderr, "✓ Migrated config to multi-context format (created 'default' context)\n")
	}

	// Ensure current context is set
	if config.CurrentContext == "" && len(config.Contexts) > 0 {
		// Set first available context as current
		for name := range config.Contexts {
			config.CurrentContext = name
			break
		}
	}

	return nil
}

// ConfigExists checks if a configuration file already exists at the given path
func ConfigExists(path string) (bool, string, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return false, "", err
		}
		path = filepath.Join(home, ".orthanc-cli.yaml")
	}

	if _, err := os.Stat(path); err == nil {
		return true, path, nil
	}
	return false, path, nil
}

// SaveConfig creates a default configuration file
func SaveConfig(path string) error {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, ".orthanc-cli.yaml")
	}

	config := `# Orthanc CLI Configuration
# Contexts allow you to manage multiple Orthanc server configurations
contexts:
  local:
    orthanc:
      url: "http://localhost:8042"
      username: "orthanc"
      password: "orthanc"
      insecure: false

# The currently active context
current-context: local

# Output configuration
output:
  json: false  # Set to true to output all results in JSON format by default
`

	return os.WriteFile(path, []byte(config), 0644)
}

// SaveConfigToFile saves the current config to file
func SaveConfigToFile(config *Config, path string) error {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, ".orthanc-cli.yaml")
	}

	viper.Set("contexts", config.Contexts)
	viper.Set("current-context", config.CurrentContext)
	viper.Set("output", config.Output)
	viper.Set("orthanc", nil) // Clear legacy field

	return viper.WriteConfigAs(path)
}
