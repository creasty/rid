package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/k0kubun/pp"

	"github.com/creasty/rid/docker"
	"github.com/creasty/rid/project"
	"github.com/creasty/rid/util"
)

const helpTemplate = `Execute commands via docker-compose

Usage:
    {{ .Name }} COMMAND [args...]
    {{ .Name }} COMMAND -h|--help
    {{ .Name }} [options]

Options:
    -h, --help     Show this
    -v, --version  Show {{ .Name }} version
        --debug    Debug context and configuration

Commands:
{{- range $name, $sub := .Command }}
    {{ printf $.NameFormat $name }}{{ if ne $sub.Summary "" }} # {{ $sub.Summary }}{{ end }}
{{- end }}
`

// CLI is an object holding states
type CLI struct {
	Context        *project.Context
	Config         *project.Config
	Args           []string
	Envs           []string
	RunInContainer bool

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewCLI creates a new CLI instance
func NewCLI(ctx *project.Context, cfg *project.Config, args []string) *CLI {
	return &CLI{
		Context:        ctx,
		Config:         cfg,
		Args:           args[1:],
		Envs:           make([]string, 0),
		RunInContainer: true,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run executes commands
func (c *CLI) Run() error {
	c.setup()
	c.parseEnvs()
	c.substituteCommand()

	switch c.Args[0] {
	case "-h", "--help", ".help":
		return c.ExecHelp()
	case "-v", "--version":
		return c.ExecVersion()
	case "--debug":
		return c.ExecDebug()
	case ".sub-help":
		return c.ExecSubHelp()
	}

	if c.RunInContainer {
		return c.runDockerExec(c.Args[0], c.Args[1:]...)
	}

	return c.run(c.Args[0], c.Args[1:]...)
}

func (c *CLI) setup() {
	os.Setenv("COMPOSE_PROJECT_NAME", c.Config.ProjectName)
	os.Setenv("DOCKER_HOST_IP", c.Context.IP)
}

func (c *CLI) parseEnvs() {
	i := 0
	for _, a := range c.Args {
		if strings.Contains(a, "=") {
			c.Envs = append(c.Envs, a)
		} else {
			break
		}
		i++
	}

	c.Args = c.Args[i:]
}

func (c *CLI) substituteCommand() {
	if len(c.Args) == 0 {
		c.Args = []string{".help"}
		return
	}

	if cmd, ok := c.Context.Command[c.Args[0]]; ok {
		c.Args[0] = cmd.Name
		c.RunInContainer = cmd.RunInContainer

		if cmd.HelpFile != "" && len(c.Args) > 1 {
			switch c.Args[1] {
			case "-h", "--help":
				c.Args = []string{".sub-help", cmd.HelpFile}
			}
		}
	}
}

func (c *CLI) run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if name == "docker-compose" {
		cmd.Dir = c.Context.BaseDir
	} else {
		cmd.Dir = c.Context.RootDir
	}
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd.Run()
}

func (c *CLI) runDockerExec(name string, args ...string) error {
	if err := c.run("docker-compose", "up", "-d", "--remove-orphans"); err != nil {
		return err
	}

	cid, err := docker.FindContainerByService(c.Config.ProjectName, c.Config.MainService, 1)
	if err != nil {
		return err
	}

	dockerArgs := []string{"exec", "-it"}
	{
		for _, e := range c.Envs {
			dockerArgs = append(dockerArgs, "-e", e)
		}

		dockerArgs = append(dockerArgs, cid)
	}

	args = append([]string{name}, args...)
	args = append(dockerArgs, args...)

	return c.run("docker", args...)
}

// ExecVersion prints version info
func (c *CLI) ExecVersion() error {
	fmt.Fprintf(c.Stdout, "%s (revision %s)\n", Version, Revision)
	return nil
}

// ExecDebug prints internal state objects
func (c *CLI) ExecDebug() error {
	pp.Fprintln(c.Stdout, c.Context)
	pp.Fprintln(c.Stdout, c.Config)
	return nil
}

// ExecHelp shows help contents
func (c *CLI) ExecHelp() error {
	maxNameLen := 0
	for name := range c.Context.Command {
		if l := len(name); l > maxNameLen {
			maxNameLen = l
		}
	}

	for _, cmd := range c.Context.Command {
		if cmd.HelpFile == "" {
			continue
		}
		cmd.Summary, _ = util.LoadHelpFile(cmd.HelpFile)
	}

	tmpl := template.Must(template.New("help").Parse(helpTemplate))
	return tmpl.Execute(c.Stderr, map[string]interface{}{
		"Command":    c.Context.Command,
		"NameFormat": fmt.Sprintf("%%-%ds", maxNameLen+1),
		"Name":       "rid",
	})
}

// ExecSubHelp shows help contents for a custom sub-command
func (c *CLI) ExecSubHelp() error {
	_, description := util.LoadHelpFile(c.Args[1])
	fmt.Fprint(c.Stderr, description)
	return nil
}
