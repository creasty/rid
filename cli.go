package main

import (
	"os"
	"os/exec"

	"github.com/k0kubun/pp"
)

type CLI struct {
	*Context
	Args []string
}

func NewCLI(ctx *Context, args []string) *CLI {
	return &CLI{
		Context: ctx,
		Args:    args[1:],
	}
}

func (c *CLI) Run() error {
	c.setup()

	if ok, err := c.ExecHelp(); ok || err != nil {
		return err
	}

	c.SubstituteCommand()

	cmd := exec.Command(c.Args[0], c.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *CLI) setup() {
	os.Chdir(c.BaseDir)

	os.Setenv("COMPOSE_PROJECT_NAME", c.Config.ProjectName)
	os.Setenv("DOCKER_HOST_IP", c.IP)

	pp.Println(c.Context)
	pp.Println(c.Args)
}

func (c *CLI) ExecHelp() (bool, error) {
	if len(c.Args) > 0 && c.Args[0] != "help" {
		return false, nil
	}

	println("help")
	return true, nil
}

func (c *CLI) SubstituteCommand() {
	switch c.Args[0] {
	case "compose":
		c.Args[0] = "docker-compose"
	default:
		if cmd, ok := c.Executable[c.Args[0]]; ok {
			c.Args[0] = cmd
		}
	}
}