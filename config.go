package main

import (
	"io/ioutil"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
)

// Default config values
const (
	DefaultMainService = "app"
	DefaultDataService = "volume"
)

// Config is a configuration object which parameters are loaded from yaml file
type Config struct {
	ProjectName string `json:"project_name" valid:"required"`
	MainService string `json:"main_service"`
	DataService string `json:"data_service"`
}

// NewConfig creates a new Config instance from a file and validates its parameters
func NewConfig(file string) (*Config, error) {
	c := &Config{}

	f, err := os.Open(file)
	if err != nil {
		return c, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return c, err
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		return c, err
	}

	if _, err := govalidator.ValidateStruct(c); err != nil {
		return c, err
	}

	c.setDefault()

	return c, nil
}

func (c *Config) setDefault() {
	if c.MainService == "" {
		c.MainService = DefaultMainService
	}
	if c.DataService == "" {
		c.DataService = DefaultDataService
	}
}
