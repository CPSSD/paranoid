package commands

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

//Init a new paranoid file system
func Init(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(0)
	}

	directory := args[0]
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatalln("FATAL :", err)
	}

	if len(files) != 0 {
		log.Fatalln("FATAL : given directory is not empty")
	}

	cmd := exec.Command("pfsm", "init", directory)
	err = cmd.Run()
	if err != nil {
		log.Fatalln("FATAL : ", err)
	}
}
