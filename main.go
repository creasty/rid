package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/creasty/rid/cli"
	"github.com/creasty/rid/project"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		exit(err)
	}

	ctx, err := project.NewContext(wd)
	if err != nil {
		exit(err)
	}

	var cfg *project.Config
	if ctx.ConfigFile != "" {
		cfg, err = project.NewConfig(ctx.ConfigFile)
		if err != nil {
			exit(err)
		}
	}

	c := cli.NewCLI(ctx, cfg, os.Args)
	if err := c.Run(); err != nil {
		exit(err)
	}
}

func exit(err error) {
	if _, ok := err.(*exec.ExitError); !ok {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
