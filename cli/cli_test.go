package cli

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/k0kubun/pp"
)

func init() {
	pp.ColoringEnabled = false
}

func TestNewCLI(t *testing.T) {
	cli := NewCLI(&Context{}, &Config{}, []string{"rid", "foo", "bar"})

	if !reflect.DeepEqual(cli.Args, []string{"foo", "bar"}) {
		t.Error("it should drop the first argument")
	}

	if cli.Stdin != os.Stdin {
		t.Error("it should initialize Stdin to the OS's standard input")
	}
	if cli.Stdout != os.Stdout {
		t.Error("it should initialize Stdin to the OS's standard output")
	}
	if cli.Stderr != os.Stderr {
		t.Error("it should initialize Stdin to the OS's standard error")
	}
}

func TestCLI_setup(t *testing.T) {
	cli := NewCLI(&Context{
		IP: "192.168.0.1",
	}, &Config{
		ProjectName: "myproject",
	}, []string{"rid"})

	reset := setTestEnvs(map[string]string{
		"COMPOSE_PROJECT_NAME": "",
		"DOCKER_HOST_IP":       "",
	})
	defer reset()

	cli.setup()
	if os.Getenv("COMPOSE_PROJECT_NAME") != cli.Config.ProjectName {
		t.Error("it should set COMPOSE_PROJECT_NAME")
	}
	if os.Getenv("DOCKER_HOST_IP") != cli.Context.IP {
		t.Error("it should set DOCKER_HOST_IP")
	}
}

func TestCLI_parseEnvs(t *testing.T) {
	cli := NewCLI(&Context{}, &Config{}, []string{"rid"})

	t.Run("no envs", func(t *testing.T) {
		cli.Args = []string{"foo", "bar"}
		cli.parseEnvs()

		if !reflect.DeepEqual(cli.Args, []string{"foo", "bar"}) {
			t.Error("it should not alternate args")
		}
	})

	t.Run("envs after command", func(t *testing.T) {
		cli.Args = []string{"foo", "bar", "AAA=123"}
		cli.parseEnvs()

		if !reflect.DeepEqual(cli.Args, []string{"foo", "bar", "AAA=123"}) {
			t.Error("it should not alternate args")
		}
	})

	t.Run("envs before command", func(t *testing.T) {
		cli.Args = []string{"AAA=123", "BBB=456", "foo", "bar"}
		cli.parseEnvs()

		if !reflect.DeepEqual(cli.Args, []string{"foo", "bar"}) {
			t.Error("it should omit env args")
		}
		if !reflect.DeepEqual(cli.Envs, []string{"AAA=123", "BBB=456"}) {
			t.Error("it should parse envs")
		}
	})
}

func TestCLI_substituteCommand(t *testing.T) {
	cli := NewCLI(&Context{
		Command: map[string]*Command{
			"host": {
				Name:           "script/host",
				RunInContainer: false,
				HelpFile:       "/path/to/help.txt",
			},
			"container": {
				Name:           "script/container",
				RunInContainer: true,
			},
		},
	}, &Config{}, []string{"rid"})

	t.Run("no args", func(t *testing.T) {
		cli.Args = []string{}
		cli.substituteCommand()

		if !reflect.DeepEqual(cli.Args, []string{".help"}) {
			t.Error("it should subsitute to .help")
		}
	})

	t.Run("host command", func(t *testing.T) {
		cli.Args = []string{"host", "foo", "bar"}
		cli.substituteCommand()

		if !reflect.DeepEqual(cli.Args, []string{"script/host", "foo", "bar"}) {
			t.Error("it should subsitute")
		}
		if cli.RunInContainer != false {
			t.Error("it should make RunInContainer false")
		}
	})

	t.Run("container command", func(t *testing.T) {
		cli.Args = []string{"container", "foo", "bar"}
		cli.substituteCommand()

		if !reflect.DeepEqual(cli.Args, []string{"script/container", "foo", "bar"}) {
			t.Error("it should subsitute")
		}
		if cli.RunInContainer != true {
			t.Error("it should make RunInContainer true")
		}
	})

	t.Run("sub-command help (no file)", func(t *testing.T) {
		cli.Args = []string{"container", "-h"}
		cli.substituteCommand()

		if !reflect.DeepEqual(cli.Args, []string{"script/container", "-h"}) {
			t.Error("it should subsitute the command alone")
		}
	})

	t.Run("sub-command help (with file)", func(t *testing.T) {
		cli.Args = []string{"host", "-h"}
		cli.substituteCommand()

		if !reflect.DeepEqual(cli.Args, []string{".sub-help", cli.Context.Command["host"].HelpFile}) {
			t.Fatal("it should subsitute to .sub-help")
		}

		cli.Args = []string{"host", "--help"}
		cli.substituteCommand()

		if !reflect.DeepEqual(cli.Args, []string{".sub-help", cli.Context.Command["host"].HelpFile}) {
			t.Fatal("it should subsitute to .sub-help")
		}
	})
}

func TestCLI_ExecVersion(t *testing.T) {
	stdout := new(bytes.Buffer)

	cli := NewCLI(&Context{}, &Config{}, []string{"rid"})
	cli.Stdout = stdout

	if err := cli.ExecVersion(); err != nil {
		t.Fatalf("it should not return error: %v", err)
	}

	if !strings.Contains(stdout.String(), Version) {
		t.Error("it should print a version")
	}
	if !strings.Contains(stdout.String(), Revision) {
		t.Error("it should print a revision")
	}
}

func TestCLI_ExecDebug(t *testing.T) {
	stdout := new(bytes.Buffer)

	cli := NewCLI(&Context{}, &Config{}, []string{"rid"})
	cli.Stdout = stdout

	if err := cli.ExecDebug(); err != nil {
		t.Fatalf("it should not return error: %v", err)
	}

	if !strings.Contains(stdout.String(), "&cli.Context{") {
		t.Error("it should dump a Context object")
	}
	if !strings.Contains(stdout.String(), "&cli.Config{") {
		t.Error("it should dump a Config object")
	}
}

func TestCLI_ExecHelp(t *testing.T) {
	stderr := new(bytes.Buffer)

	cli := NewCLI(&Context{
		Command: map[string]*Command{
			"foobar": {
				Name: "script/foobar",
			},
		},
	}, &Config{}, []string{"rid"})
	cli.Stderr = stderr

	if err := cli.ExecHelp(); err != nil {
		t.Fatalf("it should not return error: %v", err)
	}

	if !strings.Contains(stderr.String(), "foobar") {
		t.Error("it should print a help")
	}
}

func TestCLI_ExecSubHelp(t *testing.T) {
	stderr := new(bytes.Buffer)

	cli := NewCLI(&Context{}, &Config{}, []string{"rid"})
	cli.Args = []string{".sub-help", "../rid/libexec/rid-sample.txt"}
	cli.Stderr = stderr

	if err := cli.ExecHelp(); err != nil {
		t.Fatalf("it should not return error: %v", err)
	}

	if !strings.Contains(stderr.String(), "Usage:") {
		t.Error("it should print a help of sub command")
	}
}

func setTestEnvs(kv map[string]string) func() {
	original := make(map[string]string)
	for k, v := range kv {
		original[k] = os.Getenv(k)
		os.Setenv(k, v)
	}
	return func() {
		for k := range kv {
			os.Setenv(k, original[k])
		}
	}
}
