package cli

import (
	"io"

	"github.com/spf13/afero"

	"github.com/creasty/rid/pkg/app"
)

// CLI is an interface for CLI
type CLI interface {
	Run() error
}

// New creates a CLI
func New(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	workingDir string,
	aferoFs afero.Fs,
) CLI {
	return &cli{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		container: app.NewDIContainer(
			stdin,
			stdout,
			stderr,
			workingDir,
			aferoFs,
		),
	}
}

type cli struct {
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	container app.DIContainer
}
