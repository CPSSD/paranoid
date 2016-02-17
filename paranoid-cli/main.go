package main

import (
	"github.com/codegangsta/cli"
	pfscommands "github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/paranoid-cli/commands"
	"github.com/cpssd/paranoid/paranoid-cli/tls"
	"os"
)

func main() {
	pfscommands.Log = logger.New("libpfs", "paranoidcli", os.DevNull)
	commands.Log = logger.New("command", "paranoidcli", os.DevNull)
	tls.Log = logger.New("tls", "paranoidcli", os.DevNull)

	app := cli.NewApp()
	app.Name = "paranoid-cli"
	app.HelpName = "paranoid-cli"
	app.Version = "0.4.0"
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "u, unsecure",
					Usage: "disable TLS/SSL for this filesystem's network services",
				},
				cli.StringFlag{
					Name:  "cert",
					Usage: "path to existing certificate file",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "path to existing key file",
				},
				cli.StringFlag{
					Name:  "p, pool",
					Usage: "name of the pool, defaults to random",
				},
			},
		},
		{
			Name:      "mount",
			Usage:     "mount a paranoid file system",
			ArgsUsage: "pfs-name mountpoint",
			Action:    commands.Mount,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "n, noprompt",
					Usage: "disable the prompt when attempting to mount a PFS without TLS/SSL",
				},
				cli.StringFlag{
					Name:  "i, interface",
					Usage: "name a network interface over which to make connections. Defaults to default interface",
				},
				cli.StringFlag{
					Name:  "d, discovery-addr",
					Usage: "Use a custom discovery server. Specified with ip:port. Defaults to public discovery server",
				},
			},
		},
		{
			Name:      "secure",
			Usage:     "secure an unsecured paranoid file system",
			ArgsUsage: "pfs-name",
			Action:    commands.Secure,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "f, force",
					Usage: "overwrite any existing cert or key files",
				},
				cli.StringFlag{
					Name:  "cert",
					Usage: "path to existing certificate file",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "path to existing key file",
				},
			},
		},
		{
			Name:      "restart",
			Usage:     "restarts the networking services",
			ArgsUsage: "pfs-name",
			Action:    commands.Restart,
		},
		{
			Name:      "automount",
			Usage:     "automount a paranoid file system with previous settings",
			ArgsUsage: "pfs-name",
			Action:    commands.AutoMount,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "n, noprompt",
					Usage: "disable the prompt when attempting to mount a PFS without TLS/SSL",
				},
				cli.StringFlag{
					Name:  "i, interface",
					Usage: "name a network interface over which to make connections. Defaults to default interface",
				},
			},
		},
		{
			Name:      "unmount",
			ArgsUsage: "pfs-name",
			Usage:     "unmount a paranoid file system",
			Action:    commands.Unmount,
		},
		{
			Name:   "list",
			Usage:  "list all paranoid file systems",
			Action: commands.List,
		},
		{
			Name:      "delete",
			ArgsUsage: "pfs-name",
			Usage:     "delete a paranoid file system",
			Action:    commands.Delete,
		},
	}
	app.Run(os.Args)
}
