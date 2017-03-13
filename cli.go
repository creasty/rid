package main

import (
	"os"
	"os/exec"

	"github.com/k0kubun/pp"
)

type CLI struct {
	Context *Context
	Args    []string
}

func NewCLI(ctx *Context, args []string) *CLI {
	return &CLI{
		Context: ctx,
		Args:    args[1:],
	}
}

func (c *CLI) Run() error {
	c.setup()

	// cmd := exec.Command("docker-compose", "-h")
	cmd := exec.Command(c.Args[0], c.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *CLI) setup() {
	os.Chdir(c.Context.BaseDir)

	os.Setenv("COMPOSE_PROJECT_NAME", c.Context.Config.ProjectName)
	os.Setenv("DOCKER_HOST_IP", c.Context.IP)

	pp.Println(c.Context)
	pp.Println(c.Args)
}
