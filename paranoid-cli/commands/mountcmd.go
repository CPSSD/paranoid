package commands

import (
	"bufio"
	"fmt"
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
	pidPath := path.Join(usr.HomeDir, ".pfs", args[0], "meta", "pfsd.pid")
	if _, err := os.Stat(pidPath); err != nil { // if PFSD is killed -9 then the pid file will still exist.
		log.Println("INFO: outdated PID file exists, Removing")
		err = os.Remove(pidPath)
		if err != nil {
			log.Println("INFO: PID file cannot be removed")
		}
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	pfsDir := path.Join(usr.HomeDir, ".pfs", args[1])

	if _, err := os.Stat(path.Join(pfsDir, "contents")); os.IsNotExist(err) {
		log.Fatalln("FATAL : PFS directory does not include contents directory")
	}
	if _, err := os.Stat(path.Join(pfsDir, "meta")); os.IsNotExist(err) {
		log.Fatalln("FATAL : PFS directory does not include meta directory")
	}
	if _, err := os.Stat(path.Join(pfsDir, "names")); os.IsNotExist(err) {
		log.Fatalln("FATAL : PFS directory does not include names directory")
	}
	if _, err := os.Stat(path.Join(pfsDir, "inodes")); os.IsNotExist(err) {
		log.Fatalln("FATAL : PFS directory does not include inodes directory")
	}

	splitAddress := strings.Split(serverAddress, ":")
	if len(splitAddress) != 2 {
		log.Fatalln("FATAL : server-address in wrong format")
	}

	cmd := exec.Command("pfsm", "mount", pfsDir, splitAddress[0], splitAddress[1], args[2])
	err = cmd.Run()
	if err != nil {
		log.Fatalln("FATAL error running pfsm mount command : ", err)
	}

	outfile, err := os.Create(path.Join(pfsDir, "meta", "logs", "pfsdLog.txt"))

	if err != nil {
		log.Fatalln("FATAL error creating output file")
	}

	if !c.GlobalBool("networkoff") {
		// Check if the cert and key files are present.
		certPath := path.Join(pfsDir, "meta", "cert.pem")
		keyPath := path.Join(pfsDir, "meta", "key.pem")
		if pathExists(certPath) && pathExists(keyPath) {
			log.Println("INFO: Starting PFSD in secure mode.")
			//TODO(terry): Add a way to check if the given cert is its own CA,
			// and skip validation based on that.
			cmd = exec.Command("pfsd", "-cert", certPath, "-key", keyPath,
				"-skip_verification", pfsDir, splitAddress[0], splitAddress[1])
			cmd.Stderr = outfile
			err = cmd.Start()
		} else {
			// Start in unsecure mode
			if !c.Bool("noprompt") {
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Print("Starting networking in unsecure mode. Are you sure? [y/N] ")
				scanner.Scan()
				answer := strings.ToLower(scanner.Text())
				if !strings.HasPrefix(answer, "y") {
					fmt.Println("Okay. Exiting ...")
					os.Exit(1)
				}
			}
			log.Println("INFO: Starting PFSD in unsecure mode.")
			cmd = exec.Command("pfsd", pfsDir, splitAddress[0], splitAddress[1])
			cmd.Stderr = outfile
			err = cmd.Start()
		}
	}

	pfiArgs := []string{pfsDir, args[2]}
	var pfiFlags []string
	if c.GlobalBool("verbose") {
		pfiFlags = append(pfiFlags, "-v")
	}
	if c.GlobalBool("networkoff") {
		pfiFlags = append(pfiFlags, "-n")
	}

	cmd = exec.Command("pfi", append(pfiFlags, pfiArgs...)...)
	outfile, err = os.Create(path.Join(pfsDir, "meta", "logs", "pfiLog.txt"))

	if err != nil {
		log.Fatalln("FATAL error creating output file")
	}
	cmd.Stderr = outfile
	err = cmd.Start()
	if err != nil {
		log.Fatalln("FATAL error running pfi command : ", err)
	}
}
