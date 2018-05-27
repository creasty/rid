package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/creasty/rid/pkg/app"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		exit(err)
	}

	c := cli.New(
		os.Stdin,
		os.Stdout,
		os.Stderr,
		wd,
		afero.NewOsFs(),
	)
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
