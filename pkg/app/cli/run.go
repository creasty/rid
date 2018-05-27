package cli

import (
	"github.com/k0kubun/pp"
)

func (c *cli) Run(args []string) error {
	c.ParseArgs(args)

	pp.Println(args)
	pp.Println(c.Config)
	c.RunUsecase.Exec(args[0], args[1:]...)
	return nil
}
