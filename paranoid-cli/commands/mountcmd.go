package commands

import (
	"github.com/codegangsta/cli"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

//Talks to the necessary other programs to mount the pfs filesystem.
//If the file system doesn't exist it creates it.
func Mount(c *cli.Context) {
	args := c.Args()
	if len(args) < 4 {
		cli.ShowAppHelp(c)
		os.Exit(0)
	}

	serverAddress := args[1]
	_, err := net.DialTimeout("tcp", serverAddress, time.Duration(5*time.Second))
	if err != nil {
		log.Fatalln("FATAL : unable to reach server", err)
	}

	directory := args[2]
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

	cmd := exec.Command("pfsm", "mount", directory, splits[0], splits[1])
	err = cmd.Run()
	if err != nil {
		log.Fatalln("FATAL error running pfsm mount command : ", err)
	}

	port, err := strconv.Atoi(args[0])
	if err != nil || port < 1 || port > 65535 {
		log.Fatalln("FATAL: port must be a number between 1 and 65535, inclusive.")
	}

	outfile, err := os.Create(path.Join(directory, "meta", "logs", "pfsdLog.txt"))

	if err != nil {
		log.Fatalln("FATAL error creating output file")
	}

	cmd = exec.Command("pfsd", args[0], directory, splits[0], splits[1])
	cmd.Stderr = outfile
	err = cmd.Start()

	cmd = exec.Command("pfi", directory, args[3])
	if c.GlobalBool("verbose") {
		cmd = exec.Command("pfi", "-v", directory, args[3])
	}
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
