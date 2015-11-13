package util

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"github.com/hanwen/go-fuse/fuse"
	"log"
)

var LogOutput bool
var MountPoint string
var PfsDirectory string

//LogMessage checks if the -v flag was specified and either logs or doesnt log the message
func LogMessage(message string) {
	if LogOutput {
		log.Println(message)
	}
}

func GetFuseReturnCode(retcode int) fuse.Status {
	switch retcode {
	case returncodes.ENOENT:
		return fuse.ENOENT
	case returncodes.EACCES:
		return fuse.EACCES
	case returncodes.EEXIST:
		return fuse.EIO
	default:
		return fuse.OK
	}
}
