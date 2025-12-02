package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// These will be set by the main package
	version   string
	commit    string
	buildTime string
)

// SetVersionInfo sets the version information from main package
func SetVersionInfo(v, c, b string) {
	version = v
	commit = c
	buildTime = b
}

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display the version, commit, and build time of the Orthanc CLI tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Orthanc CLI\n")
			fmt.Printf("  Version:    %s\n", version)
			fmt.Printf("  Commit:     %s\n", commit)
			fmt.Printf("  Build Time: %s\n", buildTime)
		},
	}

	return cmd
}
