package commands

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
)

//Init a new paranoid file system
func Init(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(0)
	}

	pfsname := args[0]

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir := usr.HomeDir

	if _, err := os.Stat(path.Join(homeDir, "pfs")); os.IsNotExist(err) {
		err = os.Mkdir(path.Join(homeDir, "pfs"), 0777)
		if err != nil {
			log.Fatalln("FATAL : Error making pfs directory")
		}
	}

	directory := path.Join(homeDir, "pfs", pfsname)
	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		log.Fatalln("FATAL : a paranoid file system with that name already exists")
	}
	err = os.Mkdir(directory, 0777)
	if err != nil {
		log.Fatalln("FATAL : Error making pfs directory")
	}

	cmd := exec.Command("pfsm", "init", directory)
	err = cmd.Run()
	if err != nil {
		log.Fatalln("FATAL : ", err)
	}
}
