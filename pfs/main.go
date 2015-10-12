package main

import (
	"github.com/cpssd/paranoid/pfs/initcmd"
	"github.com/cpssd/paranoid/pfs/mountcmd"
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
		if onlyArgs[0] == "init" {
			initcmd.InitCommand(onlyArgs[1:])
		} else if onlyArgs[0] == "mount" {
			mountcmd.MountCommand(onlyArgs[1:])
		}
	}
}
