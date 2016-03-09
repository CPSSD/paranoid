package commands

import (
	"fmt"
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
		fmt.Println("FATAL: Error Getting Current User")
		Log.Fatal("Cannot get curent User:", err)
	}

	pidPath := path.Join(usr.HomeDir, ".pfs", args[0], "meta", "pfsd.pid")
	if _, err := os.Stat(pidPath); err == nil {
		pidByte, err := ioutil.ReadFile(pidPath)
		if err != nil {
			fmt.Println("FATAL: Can't read pid file")
			Log.Fatal("Can't read pid file", err)
		}
		pid, err := strconv.Atoi(string(pidByte))
		err = syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			fmt.Println("FATAL: Can not kill PFSD")
			Log.Fatal("Can not kill PFSD,", err)
		}
	} else {
		fmt.Println("FATAL: Could not read pid file")
		Log.Fatal("Could not read pid file:", err)
	}
}
