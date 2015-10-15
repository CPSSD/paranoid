package pfsInterface

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

/*
Stat -

description :
    Called when the attributes of a file or directory are needed.

paramenters :
    mountDir - The root directory of the file system.
    pfsLocation - The path to the pfs executable.
    name - The name of the file whos attributes are needed.

return :
    Should return a structure of the attributes, waiting on
    confirmation from the pfs team.
*/
func Stat(mountDir string, pfsLocation, name string) { // TODO: return structure
	args := fmt.Sprintf("-f stat %s %s", mountDir, name)
	command := exec.Command(pfsLocation, args)

	output, err := command.Output()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(string(output))
	// TODO: return return structure object
}
