package app

import (
	"io"
	"sync"

	"github.com/spf13/afero"

	"github.com/creasty/rid/pkg/data/repository"
	"github.com/creasty/rid/pkg/domain/model"
	"github.com/creasty/rid/pkg/domain/usecase"
	"github.com/creasty/rid/pkg/infra/docker"
	"github.com/creasty/rid/pkg/infra/fs"
)

// DIContainer instantiates and configures objects by resolving dependencies.
type DIContainer interface {
	FileSystem() fs.FileSystem
	Docker() docker.Docker
	ConfigRepository() repository.ConfigRepository
	Config() *model.Config
	RunUsecase() usecase.RunUsecase
}

// NewDIContainer creates a container
func NewDIContainer(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	workingDir string,
	aferoFs afero.Fs,
) DIContainer {
	return &diContainer{
		Stdin:      stdin,
		Stdout:     stdout,
		Stderr:     stderr,
		workingDir: workingDir,
		AferoFs:    aferoFs,
	}
}

type diContainer struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	workingDir string

	AferoFs afero.Fs

	fileSystemHolder fs.FileSystem
	fileSystemOnce   sync.Once

	configHolder *model.Config
	configOnce   sync.Once
}

func (c *diContainer) FileSystem() fs.FileSystem {
	c.fileSystemOnce.Do(func() {
		c.fileSystemHolder = fs.New(c.AferoFs)
	})
	return c.fileSystemHolder
}

func (c *diContainer) Docker() docker.Docker {
	return docker.New(
		c.Stdin,
		c.Stdout,
		c.Stderr,
		c.Config().RootDir,
		c.Config().RidDir,
	)
}

func (c *diContainer) ConfigRepository() repository.ConfigRepository {
	return repository.NewConfigRepository(c.FileSystem(), c.workingDir)
}

func (c *diContainer) Config() *model.Config {
	c.configOnce.Do(func() {
		repo := c.ConfigRepository()
		if config, err := repo.Get(); err == nil {
			c.configHolder = config
		} else {
			panic(err)
		}
	})
	return c.configHolder
}

func (c *diContainer) RunUsecase() usecase.RunUsecase {
	return usecase.NewRunUsecase(c.Config(), c.Docker())
}
