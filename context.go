package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	configFileName = "dor/dor.yml"
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
	RootDir      string
	BaseDir      string
	ConfigFile   string
	Substitution map[string]*Substitution
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
			c.RootDir = filepath.Dir(c.BaseDir)
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
			name, wrapper := removeWrapperPrefix(basename)
			if !wrapper {
				f, _ = filepath.Rel(c.RootDir, f)
			}
			c.Substitution[name] = &Substitution{
				Command:        f,
				RunInContainer: !wrapper,
			}
			continue
		}

		if strings.HasSuffix(f, ".txt") {
			help[strings.TrimSuffix(basename, ".txt")] = f
		}
	}

	for name, file := range help {
		name, _ = removeWrapperPrefix(name)
		if e, ok := c.Substitution[name]; ok {
			e.HelpFile = file
		}
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
