package usecase

import (
	"github.com/creasty/rid/pkg/domain/model"
	"github.com/creasty/rid/pkg/infra/docker"
)

// RunUsecase ...
type RunUsecase interface {
	Run() error
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

func (u *runUsecase) Run() error {
	return nil
}
