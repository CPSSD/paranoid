package commands

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
)

//Unmounts a paranoid file system
func Unmount(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(0)
	}

	cmd := exec.Command("fusermount", "-u", "-z", args[0])
	err := cmd.Run()
	if err != nil {
		log.Fatalln("FATAL : unmount failed ", err)
	}
}
