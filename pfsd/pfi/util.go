package pfi

import (
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/logger"
	"github.com/hanwen/go-fuse/fuse"
	"syscall"
)

var (
	SendOverNetwork bool
	Log             *logger.ParanoidLogger
)

func GetFuseReturnCode(retcode returncodes.Code) fuse.Status {
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
	case returncodes.EBUSY:
		return fuse.EBUSY
	default:
		return fuse.OK
	}
}
