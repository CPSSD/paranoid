package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/paranoid-cli/tls"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

func Secure(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowCommandHelp(c, "secure")
		os.Exit(0)
	}
	pfsname := args[0]
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir := usr.HomeDir
	pfsDir, err := filepath.Abs(path.Join(homeDir, ".pfs", pfsname))
	if err != nil {
		log.Fatalln("FATAL: Could not get absolute path to paranoid filesystem.")
	}
	if !pathExists(pfsDir) {
		log.Fatalln("FATAL: Paranoid filesystem does not exist:", pfsname)
	}

	certPath := path.Join(pfsDir, "meta", "cert.pem")
	keyPath := path.Join(pfsDir, "meta", "key.pem")
	if c.Bool("force") {
		os.Remove(certPath)
		os.Remove(keyPath)
	} else {
		if pathExists(certPath) || pathExists(keyPath) {
			log.Fatalln("FATAL: Paranoid filesystem already secured.",
				"Run with --force to overwrite existing security files.")
		}
	}

	if (c.String("cert") != "") && (c.String("key") != "") {
		log.Println("INFO: Using existing certificate.")
		err = os.Link(c.String("cert"), certPath)
		if err != nil {
			log.Fatalln("FATAL: Failed to copy cert file:", err)
		}
		err = os.Link(c.String("key"), keyPath)
		if err != nil {
			log.Fatalln("FATAL: Failed to copy key file:", err)
		}
	} else {
		log.Println("INFO: Generating certificate.")
		fmt.Println("Generating TLS certificate. Please follow the given instructions.")
		tls.GenCertificate(pfsDir)
	}
}
