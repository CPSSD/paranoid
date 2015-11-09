package commands

import (
	"github.com/codegangsta/cli"
	"os"
)

//Init a new paranoid file system
func Init(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(0)
	}
}
