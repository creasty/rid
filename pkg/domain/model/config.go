package model

// Config holds information of the current project
type Config struct {
	// ProjectName is used for docker-compose in order to distinguish projects in other locations.
	ProjectName string

	// MainService is a service name, in which container commands are executed.
	MainService string

	// ComposeYaml is modified content of a docker-compose yaml file.
	ComposeYaml []byte
}
