package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		exit(err)
	}

	println("hello: " + wd)
}

func exit(err error) {
	if _, ok := err.(*exec.ExitError); !ok {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}