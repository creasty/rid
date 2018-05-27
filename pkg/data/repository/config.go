//go:generate mockgen -source=config.go -package repository -destination=config_mock.go

package repository

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"

	"github.com/creasty/rid/pkg/data/entity"
	"github.com/creasty/rid/pkg/domain/model"
	"github.com/creasty/rid/pkg/infra/fs"
)

// ConfigRepository is an interface for accessing Config models
type ConfigRepository interface {
	Get() (*model.Config, error)
}

// NewConfigRepository creates an instance of ConfigRepository
func NewConfigRepository(
	workingDir string,
	fileSystem fs.FileSystem,
) ConfigRepository {
	return &configRepository{
		workingDir: workingDir,
		fileSystem: fileSystem,
	}
}

type configRepository struct {
	workingDir string
	fileSystem fs.FileSystem
}

func (r *configRepository) Get() (*model.Config, error) {
	rootInfo, err := r.getRootInfo()
	if err != nil {
		return nil, err
	}

	composeYaml, yamlData, err := r.readComposeFile(rootInfo.ComposeFile)
	if err != nil {
		return nil, err
	}

	c := &model.Config{}
	c.ProjectName = composeYaml.Rid.ProjectName
	c.MainService = composeYaml.Rid.MainService
	c.ComposeYaml = yamlData

	return c, nil
}

func (r *configRepository) readComposeFile(path string) (*entity.ComposeYaml, []byte, error) {
	c := &entity.ComposeYaml{}

	b, err := r.fileSystem.ReadFile(path)
	if err != nil {
		return c, b, errors.Wrap(err, "io error")
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		return c, b, errors.Wrap(err, "unmarshal error")
	}

	if _, err := govalidator.ValidateStruct(c); err != nil {
		return c, b, errors.Wrap(err, "invalid")
	}

	return c, b, nil
}

func (r *configRepository) getRootInfo() (*fs.RootInfo, error) {
	info, ok := r.fileSystem.LocateRoot(r.workingDir)
	if !ok {
		return nil, errors.New("Not a rid project (or any of the parent directories)")
	}
	return info, nil
}
