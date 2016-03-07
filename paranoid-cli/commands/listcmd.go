package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

//List lists all paranoid file systems
func List(c *cli.Context) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(path.Join(usr.HomeDir, ".pfs"))
	if err != nil {
		fmt.Println("Could not read the list of paranoid file systems.")
		Log.Fatal("Could not get list of paranoid file systems. Error :", err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
