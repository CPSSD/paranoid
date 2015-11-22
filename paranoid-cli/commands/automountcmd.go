package commands

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

func AutoMount(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		os.Exit(0)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	pfsMeta := path.Join(usr.HomeDir, "pfs", args[0], "meta")

	ip, err := ioutil.ReadFile(path.Join(pfsMeta, "ip"))
	if err != nil {
		log.Fatalln("FATAL : Could not get ip", err)
	}

	port, err := ioutil.ReadFile(path.Join(pfsMeta, "port"))
	if err != nil {
		log.Fatalln("FATAL : Could not get port", err)
	}

	mountpoint, err := ioutil.ReadFile(path.Join(pfsMeta, "mountpoint"))
	if err != nil {
		log.Fatalln("FATAL : Could not get mountpoint", err)
	}

	mountArgs := []string{string(ip) + string(port), args[0], string(mountpoint)}
	doMount(c, mountArgs)
}
