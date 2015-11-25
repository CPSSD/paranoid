package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os/user"
	"path"
)

//List lists all paranoid file systems
func List(c *cli.Context) {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	files, err := ioutil.ReadDir(path.Join(usr.HomeDir, ".pfs"))
	if err != nil {
		log.Fatalln("FATAL : could not get list of paranoid file systems. Error :", err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
