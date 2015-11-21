package main

import (
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/paranoid-cli/commands"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "paranoid-cli"
	app.HelpName = "paranoid-cli"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable verbose loging",
		},
		cli.BoolFlag{
			Name:  "networkoff",
			Usage: "turn off networking",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "init",
			ArgsUsage: "pfs-name",
			Usage:     "init a new paranoid file system",
			Action:    commands.Init,
		},
		{
			Name:      "mount",
			Usage:     "mount a paranoid file system",
			ArgsUsage: "discovery-server-address pfs-name mountpoint",
			Action:    commands.Mount,
		},
		{
			Name:      "automount",
			Usage:     "automount a paranoid file system with previous settings",
			ArgsUsage: "pfs-name",
			Action:    commands.AutoMount,
		},
		{
			Name:      "unmount",
			ArgsUsage: "pfs-name",
			Usage:     "unmount a paranoid file system",
			Action:    commands.Unmount,
		},
	}
	app.Run(os.Args)
}
