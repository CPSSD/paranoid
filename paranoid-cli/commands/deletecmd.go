package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

//Delete deletes a paranoid file system
func Delete(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		cli.ShowCommandHelp(c, "delete")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("FATAL: Error Getting Current User")
		Log.Fatal("Cannot get curent User:", err)
	}

	pfspath, err := filepath.Abs(path.Join(usr.HomeDir, ".pfs", args[0]))
	if err != nil {
		fmt.Println("FATAL: Paranoid file system name is incorrectly formatted")
		Log.Fatal("Given pfs-name is in incorrect format. Error : ", err)
	}
	if path.Base(pfspath) != args[0] {
		fmt.Println("FATAL: Paranoid file system name is incorrectly formatted")
		Log.Fatal("Given pfs-name is in incorrect format")
	}

	err = os.RemoveAll(path.Join(usr.HomeDir, ".pfs", args[0]))
	if err != nil {
		fmt.Println("FATAL: Could not delete given paranoid file system.")
		Log.Fatal("Could not delete given paranoid file system. Error :", err)
	}
}
