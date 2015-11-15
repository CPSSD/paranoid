package network

import (
	"github.com/cpssd/paranoid/ic/icserver"
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"github.com/cpssd/paranoid/pfsd/globals"
	"log"
)

func SendRequest(socket icserver.FileSystemMessage) {
	ips := globals.Nodes.GetAll()
	switch socket.Command {
	case "chmod":
		chmod(ips, socket.Args[0], socket.Args[1])
	case "creat":
		creat(ips, socket.Args[0], socket.Args[1])
	case "link":
		link(ips, socket.Args[0], socket.Args[1])
	case "ping":
		ping(ips)
	case "rename":
		rename(ips, socket.Args[0], socket.Args[1])
	case "truncate":
		truncate(ips, socket.Args[0], socket.Args[1])
	case "unlink":
		unlink(ips, socket.Args[0])
	case "utimes":
		utimes(ips,
			socket.Args[0],
			socket.Args[1],
			socket.Args[2],
			socket.Args[3],
			socket.Args[4])
	case "write":
		write(ips, socket.Args[0], socket.Data, socket.Args[1], socket.Args[2])
	}
}
