package cli

import (
	"github.com/k0kubun/pp"
)

func (c *cli) Run() error {
	c.RunUsecase.Run()
	pp.Println(c.Config)
	return nil
}
