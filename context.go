package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
)

const (
	configFileName = "dev/devc.yml"
	libexecDirName = "libexec"
)

type Context struct {
	BaseDir    string
	ConfigFile string
	Executable map[string]string
	Help       map[string]string
	Config     Config
	IP         string
}

func NewContext(path string) (*Context, error) {
	c := &Context{
		Executable: make(map[string]string),
		Help:       make(map[string]string),
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
	if err := c.findExecutables(); err != nil {
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

func (c *Context) findExecutables() error {
	files, err := filepath.Glob(filepath.Join(c.BaseDir, libexecDirName, "*"))
	if err != nil {
		return err
	}

	for _, f := range files {
		basename := filepath.Base(f)

		if s, err := os.Stat(f); err == nil && (s.Mode()&0111) != 0 {
			c.Executable[basename] = f
			continue
		}

		if strings.HasSuffix(f, ".txt") {
			c.Help[strings.TrimSuffix(basename, ".txt")] = f
		}
	}

	for k := range c.Help {
		if _, ok := c.Executable[k]; !ok {
			delete(c.Help, k)
		}
	}

	return nil
}

func (c *Context) loadConfig() error {
	file, err := os.Open(c.ConfigFile)
	if err != nil {
		return nil
	}

	b, err := ioutil.ReadAll(file)
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
	for _, i := range []string{"en0", "en1", "en2"} {
		cmd := exec.Command("ipconfig", "getifaddr", i)
		b, err := cmd.Output()
		if err != nil {
			continue
		}

		if len(b) > 0 {
			c.IP = strings.Trim(string(b[:]), "\n")
			return nil
		}
	}

	return errors.New("Failed to get a local IP address")
}
