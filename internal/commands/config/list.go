package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListCommand creates the config list command
func NewListCommand() *cobra.Command {
	var showPassword bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		Long:  `Display all current configuration values.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := viper.ConfigFileUsed()
			if configFile != "" {
				fmt.Printf("Configuration file: %s\n\n", configFile)
			} else {
				fmt.Println("No configuration file loaded")
				fmt.Println()
			}

			fmt.Println("Orthanc Configuration:")
			fmt.Println("----------------------")

			url := viper.GetString("orthanc.url")
			username := viper.GetString("orthanc.username")
			password := viper.GetString("orthanc.password")
			insecure := viper.GetBool("orthanc.insecure")
			jsonOutput := viper.GetBool("output.json")

			if url != "" {
				fmt.Printf("  URL:      %s\n", url)
			} else {
				fmt.Println("  URL:      (not set)")
			}

			if username != "" {
				fmt.Printf("  Username: %s\n", username)
			} else {
				fmt.Println("  Username: (not set)")
			}

			if password != "" {
				if showPassword {
					fmt.Printf("  Password: %s\n", password)
				} else {
					fmt.Println("  Password: ********")
				}
			} else {
				fmt.Println("  Password: (not set)")
			}

			fmt.Printf("  Insecure: %v\n", insecure)

			fmt.Println()
			fmt.Println("Output Configuration:")
			fmt.Println("---------------------")
			fmt.Printf("  JSON:     %v\n", jsonOutput)

			return nil
		},
	}

	cmd.Flags().BoolVar(&showPassword, "show-password", false, "Show password in plain text")

	return cmd
}
