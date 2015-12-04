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
		log.Fatal(err)
	}
	pidPath := path.Join(usr.HomeDir, ".pfs", args[0], "meta", "pfsd.pid")
	if _, err := os.Stat(pidPath); err == nil {
		pidByte, err := ioutil.ReadFile(pidPath)
		if err != nil {
			log.Fatalln("FATAL: Can't read pid file", err)
		}
		pid, err := strconv.Atoi(string(pidByte))
		err = syscall.Kill(pid, syscall.SIGHUP)
		if err != nil {
			log.Fatalln("FATAL: Can not restart PFSD,", err)
		}
	}
}
