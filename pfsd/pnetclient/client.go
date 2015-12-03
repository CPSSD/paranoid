package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/icserver"
)

func SendRequest(socket icserver.FileSystemMessage) {
	ips := globals.Nodes.GetAll()
	switch socket.Command {
	case "chmod":
		Chmod(ips, socket.Args[0], socket.Args[1])
	case "creat":
		Creat(ips, socket.Args[0], socket.Args[1])
	case "link":
		Link(ips, socket.Args[0], socket.Args[1])
	case "symlink":
		Symlink(ips, socket.Args[0], socket.Args[1])
	case "ping":
		Ping(ips)
	case "rename":
		Rename(ips, socket.Args[0], socket.Args[1])
	case "truncate":
		Truncate(ips, socket.Args[0], socket.Args[1])
	case "unlink":
		Unlink(ips, socket.Args[0])
	case "utimes":
		Utimes(ips, socket.Args[0], socket.Data)
	case "write":
		Write(ips, socket.Args[0], socket.Data, socket.Args[1], socket.Args[2])
	case "mkdir":
		Mkdir(ips, socket.Args[0], socket.Args[1])
	case "rmdir":
		Rmdir(ips, socket.Args[0])
	}
}
