package commands

import (
	"github.com/codegangsta/cli"
	"log"
	"net"
	"os"
	"os/exec"
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
		os.Exit(0)
	}

	serverAddress := args[0]
	_, err := net.DialTimeout("tcp", serverAddress, time.Duration(5*time.Second))
	if err != nil {
		log.Fatalln("FATAL : unable to reach server", err)
	}

	directory := args[1]
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

	cmd = exec.Command("pfs-network-client", "--client", directory, splits[0], splits[1])
	outfile1, err := os.Create("./out1.txt")
	if err != nil {
		log.Fatalln("FATAL error creating output file")
	}
	cmd.Stderr = outfile1
	err = cmd.Start()

	cmd = exec.Command("pfi", directory, args[2])
	outfile2, err := os.Create("./out2.txt")
	if err != nil {
		log.Fatalln("FATAL error creating output file")
	}
	cmd.Stderr = outfile2
	err = cmd.Start()
	if err != nil {
		log.Fatalln("FATAL error running pfi command : ", err)
	}
}
