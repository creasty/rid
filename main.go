package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ctx, err := NewContext(wd)
	if err != nil {
		panic(err)
	}

	cfg, err := NewConfig(ctx.ConfigFile)
	if err != nil {
		panic(err)
	}

	cli := NewCLI(ctx, cfg, os.Args)
	if err := cli.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
