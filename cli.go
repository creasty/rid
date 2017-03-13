package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/k0kubun/pp"
)

const helpTemplate = `Execute commands via docker-compose

Usage:
	devc [command]

Commands:
{{- range $name, $sub := .Substitution }}
	{{ printf "%-10s" $name }}{{ if ne $sub.Summary "" }} # {{ $sub.Summary }}{{ end }}
{{- end }}
`

type CLI struct {
	*Context
	Args           []string
	RunInContainer bool
}

func NewCLI(ctx *Context, args []string) *CLI {
	return &CLI{
		Context:        ctx,
		Args:           args[1:],
		RunInContainer: true,
	}
}

func (c *CLI) Run() error {
	c.setup()
	c.substituteCommand()

	if ok, err := c.ExecHelp(); ok || err != nil {
		return err
	}

	if c.RunInContainer {
		if err := c.exec("docker-compose", "up", "-d"); err != nil {
			return err
		}
		args := append([]string{
			"exec",
			c.Config.MainService,
		}, c.Args...)
		return c.exec("docker-compose", args...)
	}

	return c.exec(c.Args[0], c.Args[1:]...)
}

func (c *CLI) setup() {
	os.Setenv("COMPOSE_PROJECT_NAME", c.Config.ProjectName)
	os.Setenv("DOCKER_HOST_IP", c.IP)

	pp.Println(c.Context)
	pp.Println(c.Args)
}

func (c *CLI) substituteCommand() {
	if len(c.Args) == 0 {
		c.Args = []string{"help"}
		return
	}

	if s, ok := c.Substitution[c.Args[0]]; ok {
		c.Args[0] = s.Command
		c.RunInContainer = s.RunInContainer
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

func (c *CLI) ExecHelp() (bool, error) {
	if c.Args[0] != "help" {
		return false, nil
	}

	for _, s := range c.Substitution {
		if s.HelpFile == "" {
			continue
		}

		f, err := os.Open(s.HelpFile)
		if err != nil {
			continue
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			continue
		}

		s.Description = string(b[:])
		s.Summary = strings.SplitN(s.Description, "\n", 2)[0] // FIXME: consider other newline chars
	}

	tmpl := template.Must(template.New("name").Parse(helpTemplate))
	return true, tmpl.Execute(os.Stderr, c)
}
