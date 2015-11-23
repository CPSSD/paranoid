package commands

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/user"
	"path"
)

//Deletes a paranoid file system
func Delete(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(0)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	err = os.RemoveAll(path.Join(usr.HomeDir, ".pfs", args[0]))
	if err != nil {
		log.Fatalln("FATAL : count not delete given paranoid file system. Error :", err)
	}
}
