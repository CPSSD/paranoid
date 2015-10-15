package pfsPipe

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const pfsLocation string = "/home/mladen/Coding/gocode/bin/pfs"
const mountDir string = "/home/mladen/Coding/pp2pTesting"

func Stat(name string) { // TODO: return structure
	args := fmt.Printf("-f stat %s %s", mountDir, name)
	command := exec.Command(pfsLocation, args)

	output, err := command.Output()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(string(output))
	// TODO: return return structure object
}

func Read(name string, offset int, length int) []bytes {
	args := fmt.Printf("-f read %s %s", mountDir, name)
	if offset != nil {
		args += " " + string(offset)

		if length != nil {
			args += " " + string(length)
		}
	}

	command := exec.Command(pfsLocation, args)

	output, err := command.Output()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(output)
	return output
}

func Readdir(name string) {
	args := fmt.Printf("-f readdir %s", mountDir)
	command := exec.Command(pfsLocation, args)

	output, err := command.Output()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	outputString := string(output)
	filenames := outputString.Split("\n")
	fmt.Println(filenames)
	return filenames
}

func Creat(name string) {
	args := fmt.Printf("-f creat %s %s", mountDir, name)
	command := exec.Command(pfsLocation, args)

	err := command.Run()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func Write(name string, offset int, length int, data []byte) {
	// TODO: do some shit here
}
