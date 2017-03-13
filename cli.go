package main

import (
	"os"
	"os/exec"

	"github.com/k0kubun/pp"
)

type CLI struct {
	*Context
	Args      []string
	AutoStart bool
}

func NewCLI(ctx *Context, args []string) *CLI {
	return &CLI{
		Context:   ctx,
		Args:      args[1:],
		AutoStart: true,
	}
}

func (c *CLI) Run() error {
	c.setup()

	if ok, err := c.ExecHelp(); ok || err != nil {
		return err
	}

	c.substituteCommand()

	if c.AutoStart {
		if err := c.start(); err != nil {
			return err
		}
	}

	return c.run()
}

func (c *CLI) setup() {
	os.Setenv("COMPOSE_PROJECT_NAME", c.Config.ProjectName)
	os.Setenv("DOCKER_HOST_IP", c.IP)

	pp.Println(c.Context)
	pp.Println(c.Args)
}

func (c *CLI) substituteCommand() {
	switch c.Args[0] {
	case "compose":
		c.Args[0] = "docker-compose"
		c.AutoStart = false
	default:
		if cmd, ok := c.Executable[c.Args[0]]; ok {
			c.Args[0] = cmd
			c.AutoStart = false
		}
	}
}

func (c *CLI) exec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = c.BaseDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *CLI) run() error {
	return c.exec(c.Args[0], c.Args[1:]...)
}

func (c *CLI) start() error {
	return c.exec("docker-compose", "up", "-d")
}

func (c *CLI) ExecHelp() (bool, error) {
	if len(c.Args) > 0 && c.Args[0] != "help" {
		return false, nil
	}

	println("help")
	return true, nil
}
