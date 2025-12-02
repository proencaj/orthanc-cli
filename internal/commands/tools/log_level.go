package tools

import (
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// NewLogLevelCommand creates the tools log-level command with get and set subcommands
func NewLogLevelCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "log-level",
		Short: "Manage Orthanc server log level",
		Long:  `Get or set the current log level of the Orthanc server dynamically without restarting.`,
	}

	// Add subcommands
	command.AddCommand(NewLogLevelGetCommand())
	command.AddCommand(NewLogLevelSetCommand())

	return command
}

// NewLogLevelGetCommand creates the log-level get command
func NewLogLevelGetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "Get the current log level",
		Long:  `Retrieve the current log level of the Orthanc server.`,
		Example: `  # Get current log level
  orthanc tools log-level get`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return runLogLevelGet()
		},
	}

	return command
}

func runLogLevelGet() error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Get the log level
	level, err := client.GetLogLevel()
	fmt.Println(level)
	if err != nil {
		return fmt.Errorf("failed to get log level: %w", err)
	}

	fmt.Printf("Current log level: %s\n", level)

	// Provide context about the log level
	switch level {
	case types.LogLevelDefault:
		fmt.Println("  (Shows only WARNING and ERROR messages)")
	case types.LogLevelVerbose:
		fmt.Println("  (Includes INFO level messages)")
	case types.LogLevelTrace:
		fmt.Println("  (Includes detailed TRACE level messages for debugging)")
	}

	return nil
}

// NewLogLevelSetCommand creates the log-level set command
func NewLogLevelSetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "set <level>",
		Short: "Set the log level",
		Long: `Dynamically change the log level of the Orthanc server without restarting.

Valid log levels:
  - default: Shows only WARNING and ERROR messages
  - verbose: Adds INFO level messages
  - trace:   Includes detailed TRACE level messages for debugging

Note: This resets all category-specific log levels.`,
		Example: `  # Set log level to verbose
  orthanc tools log-level set verbose

  # Set log level to trace for debugging
  orthanc tools log-level set trace

  # Set log level back to default
  orthanc tools log-level set default`,
		Args: cobra.ExactArgs(1),
		ValidArgs: []string{"default", "verbose", "trace"},
		RunE: func(c *cobra.Command, args []string) error {
			return runLogLevelSet(args[0])
		},
	}

	return command
}

func runLogLevelSet(levelStr string) error {
	// Validate log level
	var level types.LogLevel
	switch levelStr {
	case "default":
		level = types.LogLevelDefault
	case "verbose":
		level = types.LogLevelVerbose
	case "trace":
		level = types.LogLevelTrace
	default:
		return fmt.Errorf("invalid log level '%s', must be one of: default, verbose, trace", levelStr)
	}

	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Set the log level
	err = client.SetLogLevel(level)
	if err != nil {
		return fmt.Errorf("failed to set log level: %w", err)
	}

	fmt.Printf("Log level set to: %s\n", level)

	// Provide context about the new log level
	switch level {
	case types.LogLevelDefault:
		fmt.Println("  (Shows only WARNING and ERROR messages)")
	case types.LogLevelVerbose:
		fmt.Println("  (Includes INFO level messages)")
	case types.LogLevelTrace:
		fmt.Println("  (Includes detailed TRACE level messages for debugging)")
	}

	fmt.Println("\nNote: This change is temporary and will be reset if the server restarts.")
	fmt.Println("To make it permanent, update the configuration file and restart the server.")

	return nil
}
