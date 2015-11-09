package commands

import (
	"github.com/codegangsta/cli"
	"os"
)

//Unmounts a paranoid file system
func Unmount(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(0)
	}
}
