package commands

import (
	"strings"
)

type programFlags struct {
	Network bool
	Fuse    bool
	Version bool
	Verbose bool
}

var Flags = programFlags{
	Network: false,
	Fuse:    false,
	Version: false,
	Verbose: false,
}

func ProcessFlags(toFlags []string) {
	for i := 0; i < len(toFlags); i++ {
		if strings.ToLower(toFlags[i]) == "-f" || strings.ToLower(toFlags[i]) == "--fuse" {
			Flags.Fuse = true
		} else if strings.ToLower(toFlags[i]) == "-n" || strings.ToLower(toFlags[i]) == "--network" {
			Flags.Network = true
		} else if strings.ToLower(toFlags[i]) == "--version" {
			Flags.Version = true
		} else if strings.ToLower(toFlags[i]) == "--verbose" {
			Flags.Verbose = true
		}
	}
}
