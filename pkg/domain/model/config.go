package model

import (
	"github.com/creasty/rid/pkg/data/entity"
)

// Config holds information of the current project
type Config struct {
	// RootDir is a path for the application's root directory
	RootDir string
	// RidDir is a path for the `rid/` directory that contains Docker related files
	RidDir string
	// ComposeFile is a path for a configuration file of docker-compose
	ComposeFile string

	// ProjectName is used for docker-compose in order to distinguish projects in other locations.
	ProjectName string
	// MainService is a service name, in which container commands are executed.
	MainService string

	// ComposeYaml is modified content of a docker-compose yaml file.
	ComposeYaml entity.Hash
}
