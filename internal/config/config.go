package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Orthanc OrthancConfig `mapstructure:"orthanc"`
	Output  OutputConfig  `mapstructure:"output"`
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

	return &config, nil
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
orthanc:
  url: "http://localhost:8042"
  username: "orthanc"
  password: "orthanc"
  insecure: false

# Output configuration
output:
  json: false  # Set to true to output all results in JSON format by default
`

	return os.WriteFile(path, []byte(config), 0644)
}
