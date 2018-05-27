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

	c := app.NewDIContainer(wd)
	c.RunUsecase().Run()
	println("hello: " + wd)
}

func exit(err error) {
	if _, ok := err.(*exec.ExitError); !ok {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
