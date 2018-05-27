//go:generate mockgen -source=config_repository.go -package repository -destination=config_repository_mock.go

package repository

import (
	"github.com/asaskevich/govalidator"
	"github.com/creasty/defaults"
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
	fileSystem fs.FileSystem,
	workingDir string,
) ConfigRepository {
	return &configRepository{
		FileSystem: fileSystem,
		workingDir: workingDir,
	}
}

type configRepository struct {
	FileSystem fs.FileSystem
	workingDir string
}

func (r *configRepository) Get() (*model.Config, error) {
	rootInfo, err := r.getRootInfo()
	if err != nil {
		return nil, err
	}

	composeYaml, err := r.readComposeFile(rootInfo.ComposeFile)
	if err != nil {
		return nil, err
	}

	if err := defaults.Set(composeYaml); err != nil {
		return nil, err
	}

	c := &model.Config{}
	c.RootDir = rootInfo.RootDir
	c.RidDir = rootInfo.RidDir
	c.ComposeFile = rootInfo.ComposeFile
	c.ProjectName = composeYaml.Rid.ProjectName
	c.MainService = composeYaml.Rid.MainService
	c.ComposeYaml = composeYaml.Raw

	return c, nil
}

func (r *configRepository) readComposeFile(path string) (*entity.ComposeYaml, error) {
	c := &entity.ComposeYaml{}

	b, err := r.FileSystem.ReadFile(path)
	if err != nil {
		return c, errors.Wrap(err, "io error")
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		return c, errors.Wrap(err, "unmarshal error")
	}

	raw := make(entity.Hash)
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return c, errors.Wrap(err, "unmarshal error")
	}

	delete(raw, "rid")
	c.Raw = raw

	if _, err := govalidator.ValidateStruct(c); err != nil {
		return c, errors.Wrap(err, "invalid")
	}

	return c, nil
}

func (r *configRepository) getRootInfo() (*fs.RootInfo, error) {
	info, ok := r.FileSystem.LocateRoot(r.workingDir)
	if !ok {
		return nil, errors.New("Not a rid project (or any of the parent directories)")
	}
	return info, nil
}
