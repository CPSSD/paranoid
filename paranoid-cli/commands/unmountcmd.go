package commands

import (
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
)

//Unmount unmounts a paranoid file system
func Unmount(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	mountpoint, err := ioutil.ReadFile(path.Join(usr.HomeDir, ".pfs", args[0], "meta", "mountpoint"))
	if err != nil {
		log.Fatalln("FATAL : Could not get mountpoint ", err)
	}

	cmd := exec.Command("fusermount", "-u", "-z", string(mountpoint))
	err = cmd.Run()
	if err != nil {
		log.Fatalln("FATAL : unmount failed ", err)
	}
	if dnetclient.Disconnect() != nil {
		log.Fatalln("Can't Disconnect from Discovery Server")
	}
}
