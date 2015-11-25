package commands

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

//Delete deletes a paranoid file system
func Delete(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	pfspath, err := filepath.Abs(path.Join(usr.HomeDir, ".pfs", args[0]))
	if err != nil {
		log.Fatalln("Given pfs-name is in incorrect format. Error : ", err)
	}
	if path.Base(pfspath) != args[0] {
		log.Fatalln("Given pfs-name is in incorrect format")
	}

	err = os.RemoveAll(path.Join(usr.HomeDir, ".pfs", args[0]))
	if err != nil {
		log.Fatalln("FATAL : count not delete given paranoid file system. Error :", err)
	}
}
