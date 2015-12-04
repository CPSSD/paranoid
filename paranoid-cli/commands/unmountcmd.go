package commands

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"syscall"
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

	pidPath := path.Join(usr.HomeDir, ".pfs", args[0], "meta", "pfsd.pid")
	if _, err := os.Stat(pidPath); err == nil {
		pidByte, err := ioutil.ReadFile(pidPath)
		if err != nil {
			log.Fatalln("FATAL: Can't read pid file", err)
		}
		pid, err := strconv.Atoi(string(pidByte))
		err = syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			log.Fatalln("FATAL : Can not kill PFSD,", err)
		}
	}

}
