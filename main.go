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

	c, err := NewContext(wd)
	if err != nil {
		panic(err)
	}

	cli := NewCLI(c, os.Args)
	if err := cli.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
