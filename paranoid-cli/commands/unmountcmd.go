package commands

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"syscall"
)

//Unmount unmounts a paranoid file system
func Unmount(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowCommandHelp(c, "unmount")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		Log.Fatal(err)
	}

	pidPath := path.Join(usr.HomeDir, ".pfs", args[0], "meta", "pfsd.pid")
	if _, err := os.Stat(pidPath); err == nil {
		pidByte, err := ioutil.ReadFile(pidPath)
		if err != nil {
			Log.Fatal("Can't read pid file", err)
		}
		pid, err := strconv.Atoi(string(pidByte))
		err = syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			Log.Fatal("Can not kill PFSD,", err)
		}
	} else {
		Log.Fatal("Could not read pid file:", err)
	}
}
