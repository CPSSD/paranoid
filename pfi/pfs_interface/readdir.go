package pfsinterface

import (
	"log"
	"os/exec"
	"strings"
)

//Readdir gets the contents of a directory from pfs
func Readdir(initDir, pfsLocation, name string) (fileNames []string) {
	command := exec.Command(pfsLocation, "-f", "readdir", initDir)

	output, err := command.Output()
	outputString := string(output)

	if err != nil {
		log.Fatalln(err)
	}
	if outputString == "" {
		return make([]string, 0)
	}

	filenames := strings.Split(outputString, "\n")
	return filenames
}
