package repository

import (
	"os"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/creasty/rid/pkg/data/entity"
	"github.com/creasty/rid/pkg/infra/fs"
)

//=== Context
//==============================================================================================
type configRepositoryTestContext struct {
	ctrl       *gomock.Controller
	fs         fs.FileSystem
	repo       ConfigRepository
	workingDir string
}

func newConfigRepositoryTestContext(t *testing.T) *configRepositoryTestContext {
	ctrl := gomock.NewController(t)

	fileSystem := fs.NewTest()

	workingDir := "/app"

	return &configRepositoryTestContext{
		ctrl:       ctrl,
		fs:         fileSystem,
		workingDir: workingDir,
		repo:       NewConfigRepository(workingDir, fileSystem),
	}
}

func (c *configRepositoryTestContext) prepareFiles() {
	composeFile, _ := yaml.Marshal(entity.Hash{
		"rid": entity.Hash{
			"project_name": "sample",
		},
	})

	c.fs.MkdirAll(c.workingDir+"/rid", os.ModeDir)
	c.fs.WriteFile(c.workingDir+"/rid/docker-compose.yml", composeFile, 0666)
}

//=== Test
//==============================================================================================
func Test_ConfigRepository_getRootInfo(t *testing.T) {
	t.Run("no rid directory", func(t *testing.T) {
		ctx := newConfigRepositoryTestContext(t)
		defer ctx.ctrl.Finish()

		_, err := ctx.repo.Get()
		assert.NotNil(t, err)
	})

	t.Run("rid directory", func(t *testing.T) {
		ctx := newConfigRepositoryTestContext(t)
		defer ctx.ctrl.Finish()

		ctx.prepareFiles()

		config, err := ctx.repo.Get()
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.NotNil(t, config)
		assert.Equal(t, "sample", config.ProjectName)
	})
}
