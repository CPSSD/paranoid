package commands

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
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
		log.Fatal("FATAL: Could not get current user:", err)
	}
	pidPath := path.Join(usr.HomeDir, ".pfs", args[0], "meta", "pfsd.pid")
	_, err = os.Stat(pidPath)
	if err != nil {
		log.Fatalln("FATAL: Could not access PID file:", err)
	}
	pidByte, err := ioutil.ReadFile(pidPath)
	if err != nil {
		log.Fatalln("FATAL: Can't read PID file:", err)
	}
	pid, err := strconv.Atoi(string(pidByte))
	err = syscall.Kill(pid, syscall.SIGHUP)
	if err != nil {
		log.Fatalln("FATAL: Can not restart PFSD:", err)
	}
}
