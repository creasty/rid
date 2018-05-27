//go:generate mockgen -source=docker.go -package docker -destination=docker_mock.go

package docker

// Docker is an interface for communicating with docker
type Docker interface {
	// GetDockerComposeVersion returns a version of docker-compose
	GetDockerComposeVersion() (string, error)

	// NormalizeProjectName normalizes a project name
	NormalizeProjectName(str string) string

	// FindContainerByService returns a container ID for service
	FindContainer(projectName, service string, num int) (string, error)
}

// New creates an instance of Docker
func New() Docker {
	return &docker{}
}

type docker struct {
}
