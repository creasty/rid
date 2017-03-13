package main

import (
	"os"
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
	if cli.Run() != nil {
		os.Exit(1)
	}
}
