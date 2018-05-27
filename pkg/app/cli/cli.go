package cli

import (
	"io"

	"github.com/creasty/rid/pkg/domain/model"
	"github.com/creasty/rid/pkg/domain/usecase"
)

// CLI is an interface for CLI
type CLI interface {
	Run(args []string) error
}

// New creates a CLI
func New(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	config *model.Config,
	runUsecase usecase.RunUsecase,
) CLI {
	return &cli{
		Stdin:      stdin,
		Stdout:     stdout,
		Stderr:     stderr,
		Config:     config,
		RunUsecase: runUsecase,
		args:       make([]string, 0),
		envs:       make([]string, 0),
	}
}

type cli struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	Config     *model.Config
	RunUsecase usecase.RunUsecase
	args       []string
	envs       []string
}
