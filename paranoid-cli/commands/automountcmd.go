package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

//AutoMount mounts a file system with the last used settings.
func AutoMount(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowCommandHelp(c, "automount")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pfsMeta := path.Join(usr.HomeDir, ".pfs", args[0], "meta")

	ip, err := ioutil.ReadFile(path.Join(pfsMeta, "ip"))
	if err != nil {
		fmt.Println("Could not get ip", err)
		os.Exit(1)
	}

	port, err := ioutil.ReadFile(path.Join(pfsMeta, "port"))
	if err != nil {
		fmt.Println("Could not get port", err)
		os.Exit(1)
	}

	mountpoint, err := ioutil.ReadFile(path.Join(pfsMeta, "mountpoint"))
	if err != nil {
		fmt.Println("Could not get mountpoint", err)
		os.Exit(1)
	}

	mountArgs := []string{string(ip) + ":" + string(port), args[0], string(mountpoint)}
	doMount(c, mountArgs)
}
