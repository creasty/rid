package cli

import (
	"github.com/k0kubun/pp"
)

func (c *cli) Run() error {
	u := c.container.RunUsecase()
	u.Run()
	pp.Println(c.container.Config())
}
