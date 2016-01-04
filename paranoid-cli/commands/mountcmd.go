package commands

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
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
		cli.ShowCommandHelp(c, "mount")
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
	pfsDir := path.Join(usr.HomeDir, ".pfs", args[1])

	if _, err := os.Stat(pfsDir); os.IsNotExist(err) {
		log.Fatalln("FATAL : PFS directory does not exist")
	}
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

	if pathExists(path.Join(pfsDir, "meta/", "pfsd.pid")) {
		err = os.Remove(path.Join(pfsDir, "meta/", "pfsd.pid"))
		if err != nil {
			log.Fatalln("FATAL : Could not remove old pfsd.pid")
		}
	}

	splitAddress := strings.Split(serverAddress, ":")
	if len(splitAddress) != 2 {
		log.Fatalln("FATAL : server-address in wrong format")
	}

	returncode, err := commands.MountCommand(pfsDir, splitAddress[0], splitAddress[1], args[2])
	if returncode != returncodes.OK {
		log.Fatalln("FATAL error running pfs mount command : ", err)
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
			pfsdArgs := []string{"-cert=" + certPath, "-key=" + keyPath, "-skip_verification",
				pfsDir, args[2], splitAddress[0], splitAddress[1]}
			var pfsdFlags []string
			if c.GlobalBool("verbose") {
				pfsdFlags = append(pfsdFlags, "-v")
			}
			cmd := exec.Command("pfsd", append(pfsdFlags, pfsdArgs...)...)
			cmd.Stderr = outfile
			err = cmd.Start()
			if err != nil {
				log.Fatalln("FATAL error running pfsd command :", err)
			}
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
			pfsdArgs := []string{pfsDir, args[2], splitAddress[0], splitAddress[1]}
			var pfsdFlags []string
			if c.GlobalBool("verbose") {
				pfsdFlags = append(pfsdFlags, "-v")
			}
			cmd := exec.Command("pfsd", append(pfsdFlags, pfsdArgs...)...)
			cmd.Stderr = outfile
			err = cmd.Start()
			if err != nil {
				log.Fatalln("FATAL error running pfsd :", err)
			}
		}
	} else {
		//No need to worry about security certs
		pfsdArgs := []string{pfsDir, args[2], splitAddress[0], splitAddress[1]}
		var pfsdFlags []string
		if c.GlobalBool("verbose") {
			pfsdFlags = append(pfsdFlags, "-v")
		}
		pfsdFlags = append(pfsdFlags, "--no_networking")
		cmd := exec.Command("pfsd", append(pfsdFlags, pfsdArgs...)...)
		cmd.Stderr = outfile
		err = cmd.Start()
		if err != nil {
			log.Fatalln("FATAL error running pfsd :", err)
		}
	}
}
