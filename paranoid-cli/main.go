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
			Name:  "verbose, v",
			Usage: "enable verbose loging",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "init",
			ArgsUsage: "pfs-directory",
			Usage:     "init a new paranoid file system",
			Action:    commands.Init,
		},
		{
			Name:      "mount",
			Usage:     "mount a paranoid file system",
			ArgsUsage: "server-address pfs-directory mountpoint",
			Action:    commands.Mount,
		},
		{
			Name:      "unmount",
			ArgsUsage: "mountpoint",
			Usage:     "unmount a paranoid file system",
			Action:    commands.Unmount,
		},
	}
	app.Run(os.Args)
}
