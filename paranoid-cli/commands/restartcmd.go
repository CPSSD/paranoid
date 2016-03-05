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

func Restart(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowCommandHelp(c, "restart")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Could not get current user:", err)
		os.Exit(1)
	}

	pidPath := path.Join(usr.HomeDir, ".pfs", args[0], "meta", "pfsd.pid")
	_, err = os.Stat(pidPath)
	if err != nil {
		fmt.Println("Could not access PID file:", err)
		os.Exit(1)
	}

	pidByte, err := ioutil.ReadFile(pidPath)
	if err != nil {
		fmt.Println("Can't read PID file:", err)
		os.Exit(1)
	}
	pid, err := strconv.Atoi(string(pidByte))
	err = syscall.Kill(pid, syscall.SIGHUP)
	if err != nil {
		fmt.Println("Can not restart PFSD:", err)
		os.Exit(1)
	}
}
