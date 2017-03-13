package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
)

const (
	configFileName = "dev/devc.yml"
	libexecDirName = "libexec"
)

type Substitution struct {
	Command        string
	Summary        string
	Description    string
	RunInContainer bool
	HelpFile       string
}

type Context struct {
	BaseDir      string
	ConfigFile   string
	Substitution map[string]*Substitution
	Config       Config
	IP           string
}

func NewContext(path string) (*Context, error) {
	c := &Context{
		Substitution: map[string]*Substitution{
			"compose": {
				Command: "docker-compose",
				Summary: "Execute docker-compose",
			},
		},
	}
	if err := c.findConfigFile(path); err != nil {
		return nil, err
	}
	if err := c.loadConfig(); err != nil {
		return nil, err
	}
	if err := c.getLocalIP(); err != nil {
		return nil, err
	}
	if err := c.findSubstitutions(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Context) findConfigFile(path string) error {
	for path != "." && path != "/" {
		configFile := filepath.Join(path, configFileName)
		if _, err := os.Stat(configFile); err == nil {
			c.ConfigFile = configFile
			c.BaseDir = filepath.Dir(configFile)
			return nil
		}

		path = filepath.Dir(path)
	}

	return fmt.Errorf("Unable to locate a config file: %s", configFileName)
}

func (c *Context) findSubstitutions() error {
	files, err := filepath.Glob(filepath.Join(c.BaseDir, libexecDirName, "*"))
	if err != nil {
		return err
	}

	help := make(map[string]string)

	for _, f := range files {
		basename := filepath.Base(f)

		if s, err := os.Stat(f); err == nil && (s.Mode()&0111) != 0 {
			c.Substitution[basename] = &Substitution{
				Command:        f,
				RunInContainer: false, // TODO
			}
			continue
		}

		if strings.HasSuffix(f, ".txt") {
			help[strings.TrimSuffix(basename, ".txt")] = f
		}
	}

	for name, file := range help {
		if e, ok := c.Substitution[name]; ok {
			e.HelpFile = file
		}
	}

	return nil
}

func (c *Context) loadConfig() error {
	f, err := os.Open(c.ConfigFile)
	if err != nil {
		return nil
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, &c.Config); err != nil {
		return err
	}

	if _, err := govalidator.ValidateStruct(c.Config); err != nil {
		return err
	}

	if c.Config.MainService == "" {
		c.Config.MainService = DefaultMainService
	}
	if c.Config.DataService == "" {
		c.Config.DataService = DefaultVolumeService
	}

	return nil
}

func (c *Context) getLocalIP() error {
	c.IP = getLocalIP()
	if c.IP == "" {
		return errors.New("Failed to get a local IP address")
	}

	return nil
}
