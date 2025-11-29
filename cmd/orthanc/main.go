package main

import (
	cmd "github.com/proencaj/orthanc-cli/internal/commands"
	"github.com/proencaj/orthanc-cli/internal/commands/instances"
	"github.com/proencaj/orthanc-cli/internal/commands/modalities"
	"github.com/proencaj/orthanc-cli/internal/commands/patients"
	"github.com/proencaj/orthanc-cli/internal/commands/series"
	"github.com/proencaj/orthanc-cli/internal/commands/studies"
)

func main() {
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

	// Register commands
	cmd.AddCommand(studies.NewStudiesCommand())
	cmd.AddCommand(series.NewSeriesCommand())
	cmd.AddCommand(patients.NewPatientsCommand())
	cmd.AddCommand(instances.NewInstancesCommand())
	cmd.AddCommand(modalities.NewModalitiesCommand())

	// Execute CLI
	cmd.Execute()
}
