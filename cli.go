package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/k0kubun/pp"
)

const helpTemplate = `Execute commands via docker-compose

Usage:
    devc COMMAND [args...]
    devc -h|--help

Commands:
{{- range $name, $sub := .Substitution }}
    {{ printf $.NameFormat $name }}{{ if ne $sub.Summary "" }} # {{ $sub.Summary }}{{ end }}
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

	switch c.Args[0] {
	case "--debug":
		return c.ExecDebug()
	case "-h", "--help", ".help":
		return c.ExecHelp()
	case ".sub-help":
		return c.ExecSubHelp()
	}

	if c.RunInContainer {
		return c.run()
	}

	return c.exec(c.Args[0], c.Args[1:]...)
}

func (c *CLI) setup() {
	os.Setenv("COMPOSE_PROJECT_NAME", c.Config.ProjectName)
	os.Setenv("DOCKER_HOST_IP", c.IP)
}

func (c *CLI) substituteCommand() {
	if len(c.Args) == 0 {
		c.Args = []string{".help"}
		return
	}

	if s, ok := c.Substitution[c.Args[0]]; ok {
		c.Args[0] = s.Command
		c.RunInContainer = s.RunInContainer

		if s.HelpFile != "" && len(c.Args) > 1 {
			switch c.Args[1] {
			case "-h", "--help":
				c.Args = []string{".sub-help", s.HelpFile}
			}
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
	if err := c.exec("docker-compose", "up", "-d"); err != nil {
		return err
	}

	args := append([]string{
		"exec",
		c.Config.MainService,
	}, c.Args...)

	return c.exec("docker-compose", args...)
}

func (c *CLI) ExecDebug() error {
	pp.Println(c.Context)
	return nil
}

func (c *CLI) ExecHelp() error {
	maxNameLen := 0
	for name := range c.Substitution {
		if l := len(name); l > maxNameLen {
			maxNameLen = l
		}
	}

	for _, s := range c.Substitution {
		if summary, _ := loadHelpFile(s.HelpFile); summary != "" {
			s.Summary = summary
		}
	}

	tmpl := template.Must(template.New("help").Parse(helpTemplate))
	return tmpl.Execute(os.Stderr, map[string]interface{}{
		"Substitution": c.Substitution,
		"NameFormat":   fmt.Sprintf("%%-%ds", maxNameLen+1),
	})
}

func (c *CLI) ExecSubHelp() error {
	_, description := loadHelpFile(c.Args[1])
	fmt.Fprintln(os.Stderr, description)
	return nil
}
