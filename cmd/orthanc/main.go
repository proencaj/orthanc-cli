package main

import (
	cmd "github.com/proencaj/orthanc-cli/internal/commands"
	"github.com/proencaj/orthanc-cli/internal/commands/dicomweb"
	"github.com/proencaj/orthanc-cli/internal/commands/instances"
	"github.com/proencaj/orthanc-cli/internal/commands/modalities"
	"github.com/proencaj/orthanc-cli/internal/commands/patients"
	"github.com/proencaj/orthanc-cli/internal/commands/series"
	"github.com/proencaj/orthanc-cli/internal/commands/servers"
	"github.com/proencaj/orthanc-cli/internal/commands/studies"
	"github.com/proencaj/orthanc-cli/internal/commands/system"
	"github.com/proencaj/orthanc-cli/internal/commands/tools"
	"github.com/proencaj/orthanc-cli/internal/commands/version"
)

// Version information (injected at build time via ldflags)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	// Set version information for the version command
	version.SetVersionInfo(Version, Commit, BuildTime)

	// Set up the client getter for studies command to avoid import cycle
	studies.SetClientGetter(cmd.GetClient)

	// Set up the client getter for series command to avoid import cycle
	series.SetClientGetter(cmd.GetClient)

	// Set up the client getter for patients command to avoid import cycle
	patients.SetClientGetter(cmd.GetClient)

	// Set up the client getter for instances command to avoid import cycle
	instances.SetClientGetter(cmd.GetClient)

	// Set up the client getter for modalities command to avoid import cycle
	modalities.SetClientGetter(cmd.GetClient)

	// Set up the client getter for tools command to avoid import cycle
	tools.SetClientGetter(cmd.GetClient)

	// Set up the client getter for system command to avoid import cycle
	system.SetClientGetter(cmd.GetClient)

	// Set up the client getter for dicomweb command to avoid import cycle
	dicomweb.SetClientGetter(cmd.GetClient)

	// Set up the client getter for servers command to avoid import cycle
	servers.SetClientGetter(cmd.GetClient)

	// Register commands
	cmd.AddCommand(studies.NewStudiesCommand())
	cmd.AddCommand(series.NewSeriesCommand())
	cmd.AddCommand(patients.NewPatientsCommand())
	cmd.AddCommand(instances.NewInstancesCommand())
	cmd.AddCommand(modalities.NewModalitiesCommand())
	cmd.AddCommand(servers.NewServersCommand())
	cmd.AddCommand(tools.NewToolsCommand())
	cmd.AddCommand(system.NewSystemCommand())
	cmd.AddCommand(version.NewVersionCommand())
	cmd.AddCommand(dicomweb.NewDicomwebCommand())

	// Execute CLI
	cmd.Execute()
}
