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
		fmt.Println(err)
		os.Exit(1)
	}

	pidPath := path.Join(usr.HomeDir, ".pfs", args[0], "meta", "pfsd.pid")
	if _, err := os.Stat(pidPath); err == nil {
		pidByte, err := ioutil.ReadFile(pidPath)
		if err != nil {
			fmt.Println("Can't read pid file", err)
			os.Exit(1)
		}
		pid, err := strconv.Atoi(string(pidByte))
		err = syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			fmt.Println("Can not kill PFSD,", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Could not read pid file:", err)
		os.Exit(1)
	}
}
