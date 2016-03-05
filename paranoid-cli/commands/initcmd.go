package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/paranoid-cli/tls"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

func cleanupPFS(pfsDir string) {
	err := os.RemoveAll(pfsDir)
	if err != nil {
		Log.Warn("Could not successfully clean up PFS directory.")
	}
}

//Init inits a new paranoid file system
func Init(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowCommandHelp(c, "init")
		os.Exit(0)
	}

	pfsname := args[0]

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	homeDir := usr.HomeDir

	if _, err := os.Stat(path.Join(homeDir, ".pfs")); os.IsNotExist(err) {
		err = os.Mkdir(path.Join(homeDir, ".pfs"), 0700)
		if err != nil {
			fmt.Println("Error making pfs directory")
			os.Exit(1)
		}
	}

	directory, err := filepath.Abs(path.Join(homeDir, ".pfs", pfsname))
	if err != nil {
		fmt.Println("Given pfs-name is in incorrect format. Error : ", err)
		os.Exit(1)
	}
	if path.Base(directory) != args[0] {
		fmt.Println("Given pfs-name is in incorrect format.")
		os.Exit(1)
	}

	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		fmt.Println("A paranoid file system with that name already exists")
		os.Exit(1)
	}
	err = os.Mkdir(directory, 0700)
	if err != nil {
		fmt.Println("Error making pfs directory : ", err)
		os.Exit(1)
	}

	returncode, err := commands.InitCommand(directory)
	if returncode != returncodes.OK {
		cleanupPFS(directory)
		fmt.Println("Error running pfs init : ", err)
		os.Exit(1)
	}

	// Either create a new pool name or use the one for a flag and save to meta/pool
	var pool string
	if pool = c.String("pool"); len(pool) == 0 {
		pool = getRandomName()
	}
	err = ioutil.WriteFile(path.Join(directory, "meta", "pool"), []byte(pool), 0600)
	if err != nil {
		fmt.Println("cannot save pool information:", err)
		os.Exit(1)
	}
	Log.Infof("Using pool name %s", pool)

	if c.Bool("unsecure") {
		fmt.Println("--unsecure specified. PFSD will not use TLS for its communication.")
		return
	}

	if (c.String("cert") != "") && (c.String("key") != "") {
		fmt.Println("Using existing certificate.")
		err = os.Link(c.String("cert"), path.Join(directory, "meta", "cert.pem"))
		if err != nil {
			cleanupPFS(directory)
			fmt.Println("Failed to copy cert file:", err)
		}
		err = os.Link(c.String("key"), path.Join(directory, "meta", "key.pem"))
		if err != nil {
			cleanupPFS(directory)
			fmt.Println("Failed to copy key file:", err)
		}
	} else {
		Log.Info("Generating certificate.")
		fmt.Println("Generating TLS certificate. Please follow the given instructions.")
		err = tls.GenCertificate(directory)
		if err != nil {
			cleanupPFS(directory)
			fmt.Println("Failed to generate certificate:", err)
			os.Exit(1)
		}
	}
}
