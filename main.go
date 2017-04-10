package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/creasty/rid/cli"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		exit(err)
	}

	ctx, err := cli.NewContext(wd)
	if err != nil {
		exit(err)
	}

	cfg, err := cli.NewConfig(ctx.ConfigFile)
	if err != nil {
		exit(err)
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
