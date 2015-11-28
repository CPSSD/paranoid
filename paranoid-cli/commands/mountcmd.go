package commands

import (
	"github.com/codegangsta/cli"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"
)

//Talks to the necessary other programs to mount the pfs filesystem.
//If the file system doesn't exist it creates it.
func Mount(c *cli.Context) {
	args := c.Args()
	if len(args) < 3 {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	doMount(c, args)
}

func doMount(c *cli.Context, args []string) {
	serverAddress := args[0]
	if c.GlobalBool("networkoff") == false {
		_, err := net.DialTimeout("tcp", serverAddress, time.Duration(5*time.Second))
		if err != nil {
			log.Fatalln("FATAL : unable to reach server", err)
		}
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	directory := path.Join(usr.HomeDir, ".pfs", args[1])

	if _, err := os.Stat(path.Join(directory, "contents")); os.IsNotExist(err) {
		log.Fatalln("FATAL : directory does not include contents directory")
	}
	if _, err := os.Stat(path.Join(directory, "meta")); os.IsNotExist(err) {
		log.Fatalln("FATAL : directory does not include meta directory")
	}
	if _, err := os.Stat(path.Join(directory, "names")); os.IsNotExist(err) {
		log.Fatalln("FATAL : directory does not include names directory")
	}
	if _, err := os.Stat(path.Join(directory, "inodes")); os.IsNotExist(err) {
		log.Fatalln("FATAL : directory does not include inodes directory")
	}

	splits := strings.Split(serverAddress, ":")
	if len(splits) != 2 {
		log.Fatalln("FATAL : server-address in wrong format")
	}

	cmd := exec.Command("pfsm", "mount", directory, splits[0], splits[1], args[2])
	err = cmd.Run()
	if err != nil {
		log.Fatalln("FATAL error running pfsm mount command : ", err)
	}

	outfile, err := os.Create(path.Join(directory, "meta", "logs", "pfsdLog.txt"))

	if err != nil {
		log.Fatalln("FATAL error creating output file")
	}

	if c.GlobalBool("networkoff") == false {
		cmd = exec.Command("pfsd", directory, splits[0], splits[1])
		cmd.Stderr = outfile
		err = cmd.Start()
	}

	pfiArgs := []string{directory, args[2]}
	var pfiFlags []string
	if c.GlobalBool("verbose") {
		pfiFlags = append(pfiFlags, "-v")
	}
	if c.GlobalBool("networkoff") {
		pfiFlags = append(pfiFlags, "-n")
	}

	cmd = exec.Command("pfi", append(pfiFlags, pfiArgs...)...)
	outfile, err = os.Create(path.Join(directory, "meta", "logs", "pfiLog.txt"))

	if err != nil {
		log.Fatalln("FATAL error creating output file")
	}
	cmd.Stderr = outfile
	err = cmd.Start()
	if err != nil {
		log.Fatalln("FATAL error running pfi command : ", err)
	}
}
