package util

import (
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/logger"
	"github.com/hanwen/go-fuse/fuse"
	"syscall"
)

var MountPoint string
var PfsDirectory string
var LogOutput bool
var SendOverNetwork bool
var Log *logger.ParanoidLogger

func GetFuseReturnCode(retcode int) fuse.Status {
	switch retcode {
	case returncodes.ENOENT:
		return fuse.ENOENT
	case returncodes.EACCES:
		return fuse.EACCES
	case returncodes.EEXIST:
		return fuse.Status(syscall.EEXIST)
	case returncodes.ENOTEMPTY:
		return fuse.Status(syscall.ENOTEMPTY)
	case returncodes.ENOTDIR:
		return fuse.ENOTDIR
	case returncodes.EISDIR:
		return fuse.Status(syscall.EISDIR)
	case returncodes.EIO:
		return fuse.EIO
	default:
		return fuse.OK
	}
}
