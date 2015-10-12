package main

import (
	"github.com/cpssd/paranoid/pfs/commands"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	var onlyArgs []string
	var onlyFlags []string
	for i := 0; i < len(args); i++ {
		if args[i][0] == '-' {
			onlyFlags = append(onlyFlags, args[i])
		} else {
			onlyArgs = append(onlyArgs, args[i])
		}
	}
	if len(onlyArgs) > 0 {
		switch onlyArgs[0] {
		case "init":
			commands.InitCommand(onlyArgs[1:])
		case "mount":
			commands.MountCommand(onlyArgs[1:])
		case "creat":
			commands.CreatCommand(onlyArgs[1:])
		case "write":
			commands.WriteCommand(onlyArgs[1:])
		case "read":
			commands.ReadCommand(onlyArgs[1:])
		case "readdir":
			commands.ReadDirCommand(onlyArgs[1:])
		default:
			log.Fatal("Given command not recognised")
		}
	} else {
		log.Fatal("No command given")
	}
}
