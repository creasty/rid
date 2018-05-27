package repository

import (
	"github.com/creasty/rid/pkg/domain/model"
)

// ConfigRepository ...
type ConfigRepository interface {
	Get() (*model.Config, error)
}

// NewConfigRepository ...
func NewConfigRepository() ConfigRepository {
	return &configRepository{}
}

type configRepository struct {
}

func (r *configRepository) Get() (*model.Config, error) {
	return nil, nil
}
