package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os/user"
	"path"
)

//List lists all paranoid file systems
func List(c *cli.Context) {
	usr, err := user.Current()
	if err != nil {
		Log.Fatal(err)
	}

	files, err := ioutil.ReadDir(path.Join(usr.HomeDir, ".pfs"))
	if err != nil {
		Log.Fatal("Could not get list of paranoid file systems. Error :", err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
