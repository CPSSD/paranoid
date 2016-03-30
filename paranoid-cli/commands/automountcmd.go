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
		fmt.Println("FATAL: Error Getting Current User")
		Log.Fatal("Cannot get curent User:", err)
	}
	pfsMeta := path.Join(usr.HomeDir, ".pfs", "filesystems", args[0], "meta")

	mountpoint, err := ioutil.ReadFile(path.Join(pfsMeta, "mountpoint"))
	if err != nil {
		fmt.Println("FATAL: PFSD Couldnt find FS mountpoint", err)
		Log.Fatal("Could not get mountpoint", err)
	}

	mountArgs := []string{args[0], string(mountpoint)}
	doMount(c, mountArgs)
}
