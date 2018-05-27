//go:generate mockgen -source=docker.go -package docker -destination=docker_mock.go

package docker

import (
	"io"
)

// Docker is an interface for communicating with docker
type Docker interface {
	// GetDockerComposeVersion returns a version of docker-compose.
	GetDockerComposeVersion() (string, error)

	// NormalizeProjectName normalizes a project name.
	NormalizeProjectName(str string) string

	// FindContainerByService returns a container ID for service.
	FindContainer(projectName, service string, num int) (string, error)

	// Prepare starts up containers as deamon.
	Prepare(dir string) error

	// Exec executes the given command in the specified container ID.
	// Extra environment variables can be passed in.
	Exec(cid string, envs []string, name string, args ...string) error
}

// New creates an instance of Docker
func New(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) Docker {
	return &docker{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
}

type docker struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
