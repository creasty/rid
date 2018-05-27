package app

import (
	"os"
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
	RunUsecase() usecase.RunUsecase
}

// NewDIContainer creates a container
func NewDIContainer(
	workingDir string,
) DIContainer {
	return &diContainer{
		workingDir: workingDir,
	}
}

type diContainer struct {
	workingDir string

	fileSystemHolder fs.FileSystem
	fileSystemOnce   sync.Once

	configHolder *model.Config
	configOnce   sync.Once
}

func (c *diContainer) FileSystem() fs.FileSystem {
	c.fileSystemOnce.Do(func() {
		c.fileSystemHolder = fs.New(c.aferoFs())
	})
	return c.fileSystemHolder
}

func (c *diContainer) aferoFs() afero.Fs {
	return afero.NewOsFs()
}

func (c *diContainer) Docker() docker.Docker {
	return docker.New(
		os.Stdin,
		os.Stdout,
		os.Stderr,
		c.Config().RootDir,
		c.Config().RidDir,
	)
}

func (c *diContainer) ConfigRepository() repository.ConfigRepository {
	return repository.NewConfigRepository(c.workingDir, c.FileSystem())
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
