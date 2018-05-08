package project

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/creasty/rid/util"
)

const (
	globalRidDirName = ".rid"
	configFileName   = "rid/config.yml"
	libexecDirName   = "libexec"
)

// Context represents a world where the command is executed
type Context struct {
	RootDir    string
	BaseDir    string
	ConfigFile string
	Command    map[string]*Command
	IP         string
}

// NewContext retrieves contextual information for a working directory, and creates a Context instance
func NewContext(path string) (*Context, error) {
	c := &Context{
		Command: map[string]*Command{
			"compose": {
				Name:    "docker-compose",
				Summary: "Execute docker-compose",
			},
		},
	}
	c.findConfigFile(path)
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
	usr, err := user.Current()
	if err != nil {
		return err
	}
	globalFiles, globalFileErr := filepath.Glob(filepath.Join(usr.HomeDir, globalRidDirName, libexecDirName, "*"))
	if globalFileErr != nil {
		return globalFileErr
	}
	localFiles, localFileErr := filepath.Glob(filepath.Join(c.BaseDir, libexecDirName, "*"))
	if localFileErr != nil {
		return localFileErr
	}
	files := append(globalFiles, localFiles...)

	help := make(map[string]string)

	for _, f := range files {
		basename := filepath.Base(f)

		if s, err := os.Stat(f); err == nil && (s.Mode()&0111) != 0 {
			name, wrapper := util.RemovePrefix("rid-", basename)
			if !wrapper {
				f, _ = filepath.Rel(c.RootDir, f)
			}
			c.Command[name] = &Command{
				Name:           f,
				RunInContainer: !wrapper,
			}
			continue
		}

		if strings.HasSuffix(f, ".txt") {
			help[strings.TrimSuffix(basename, ".txt")] = f
		}
	}

	for name, file := range help {
		name, _ = util.RemovePrefix("rid-", name)
		if e, ok := c.Command[name]; ok {
			e.HelpFile = file
		}
	}

	return nil
}

func (c *Context) getLocalIP() error {
	c.IP = util.GetLocalIP()
	if c.IP == "" {
		return errors.New("Failed to get a local IP address")
	}

	return nil
}
