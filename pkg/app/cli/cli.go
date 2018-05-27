package cli

import (
	"io"

	"github.com/creasty/rid/pkg/app"
)

type CLI interface {
	Run() error
}

func NewCLI(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	workingDir string,
) CLI {
	return &cli{
		Stdin:      stdin,
		Stdout:     stdout,
		Stderr:     stderr,
		workingDir: workingDir,
		container:  app.NewDIContainer(workingDir),
	}
}

type cli struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	workingDir string
	container  app.DIContainer
}
