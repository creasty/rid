package usecase

import (
	"github.com/creasty/rid/pkg/domain/model"
	"github.com/creasty/rid/pkg/infra/docker"

	"github.com/k0kubun/pp"
)

// RunUsecase ...
type RunUsecase interface {
	Exec(name string, args ...string) error
}

// NewRunUsecase ...
func NewRunUsecase(
	config *model.Config,
	docker docker.Docker,
) RunUsecase {
	return &runUsecase{
		Config: config,
		Docker: docker,
	}
}

type runUsecase struct {
	Config *model.Config
	Docker docker.Docker
}

func (u *runUsecase) Exec(name string, args ...string) error {
	// if err := u.Docker.Prepare(u.Config.RidDir); err != nil {
	// 	return err
	// }
	projectName := u.Docker.NormalizeProjectName(u.Config.ProjectName)

	cid, err := u.Docker.FindContainer(projectName, u.Config.MainService, 1)
	if err != nil {
		return err
	}

	pp.Println(cid)
	//u.Docker.Exec(cid, []string{}, name, args...)

	return nil
}
