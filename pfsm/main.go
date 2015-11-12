package main

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsm/commands"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	var onlyArgs []string
	var onlyFlags []string
	for i := 0; i < len(args); i++ {
		if len(args[i]) > 0 {
			if args[i][0] == '-' {
				onlyFlags = append(onlyFlags, args[i])
			} else {
				onlyArgs = append(onlyArgs, args[i])
			}
		}
	}
	commands.ProcessFlags(onlyFlags)
	if commands.Flags.Version {
		fmt.Println("pfsm v0.1.0")
		return
	}
	if commands.Flags.Network && commands.Flags.Fuse {
		log.Fatalln("Error, both network and fuse flags are set")
	}
	if commands.Flags.Verbose {
		if len(args) > 0 {
			givenCmd := args[0]
			for i := 1; i < len(args); i++ {
				givenCmd = givenCmd + " " + args[i]
			}
			log.Println("Given command : ", givenCmd)
		}
	}
	if len(onlyArgs) > 0 {
		switch onlyArgs[0] {
		case "init":
			commands.InitCommand(onlyArgs[1:])
		case "mount":
			commands.MountCommand(onlyArgs[1:])
		case "chmod":
			commands.ChmodCommand(onlyArgs[1:])
		case "access":
			commands.AccessCommand(onlyArgs[1:])
		case "creat":
			commands.CreatCommand(onlyArgs[1:])
		case "write":
			commands.WriteCommand(onlyArgs[1:])
		case "read":
			commands.ReadCommand(onlyArgs[1:])
		case "readdir":
			commands.ReadDirCommand(onlyArgs[1:])
		case "stat":
			commands.StatCommand(onlyArgs[1:])
		case "truncate":
			commands.TruncateCommand(onlyArgs[1:])
		case "utimes":
			commands.UtimesCommand(onlyArgs[1:])
		case "rename":
			commands.RenameCommand(onlyArgs[1:])
		case "unlink":
			commands.UnlinkCommand(onlyArgs[1:])
		case "link":
			commands.LinkCommand(onlyArgs[1:])
		default:
			log.Fatalln("Given command not recognised")
		}
	} else {
		log.Fatalln("No command given")
	}
}
