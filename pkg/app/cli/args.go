package cli

import (
	"strings"
)

func (c *cli) ParseArgs(args []string) {
	c.args = args
	c.parseEnvs()
	c.substituteArgs()
}

func (c *cli) parseEnvs() {
	i := 0
	for _, a := range c.args {
		if strings.Contains(a, "=") {
			c.envs = append(c.envs, a)
		} else {
			break
		}
		i++
	}

	c.args = c.args[i:]
}

func (c *cli) substituteArgs() {
	if len(c.args) == 0 {
		c.args = []string{"help"}
		return
	}

	switch c.args[0] {
	case "-h", "--help":
		c.args = []string{"help"}
	case "-v", "--version":
		c.args = []string{"version"}
	case "help":
		if len(c.args) > 1 {
			c.args = []string{"help-sub", c.args[1]}
		}
	}
}
