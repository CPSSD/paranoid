package file

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/pfi/util"
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

var mountPoint string

//ParanoidFile is a custom file struct with read and write functions
type ParanoidFile struct {
	Name string
	nodefs.File
}

//NewParanoidFile returns a new object of ParanoidFile
func NewParanoidFile(name string) nodefs.File {
	return &ParanoidFile{
		Name: name,
		File: nodefs.NewDefaultFile(),
	}
}

//Read reads a file and returns an array of bytes
func (f *ParanoidFile) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {
	util.Log.Info("Read called on file:", f.Name)
	code, err, data := commands.ReadCommand(util.PfsDirectory, f.Name, off, int64(len(buf)))
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running read command :", err)
	}

	if err != nil {
		util.Log.Error("Error running read command :", err)
	}

	copy(buf, data)
	if code != returncodes.OK {
		return nil, util.GetFuseReturnCode(code)
	}
	return fuse.ReadResultData(data), fuse.OK
}

//Write writes to a file
func (f *ParanoidFile) Write(content []byte, off int64) (uint32, fuse.Status) {
	util.Log.Info("Write called on file : " + f.Name)
	var (
		code         int
		err          error
		bytesWritten int
	)
	if util.SendOverNetwork {
		code, err, bytesWritten = util.RaftServer.RequestWriteCommand(f.Name, uint64(off), uint64(len(content)), content)
	} else {
		code, err, bytesWritten = commands.WriteCommand(util.PfsDirectory, f.Name, off, int64(len(content)), content)
	}

	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running write command :", err)
	}

	if err != nil {
		util.Log.Error("Error running write command :", err)
	}

	if code != returncodes.OK {
		return 0, util.GetFuseReturnCode(code)
	}

	return uint32(bytesWritten), fuse.OK
}

//Truncate is called when a file is to be reduced in length to size.
func (f *ParanoidFile) Truncate(size uint64) fuse.Status {
	util.Log.Info("Truncate called on file : " + f.Name)
	var code int
	var err error
	if util.SendOverNetwork {
		code, err = util.RaftServer.RequestTruncateCommand(f.Name, size)
	} else {
		code, err = commands.TruncateCommand(util.PfsDirectory, f.Name, int64(size))
	}

	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running truncate command :", err)
	}

	if err != nil {
		util.Log.Error("Error running truncate command :", err)
	}

	return util.GetFuseReturnCode(code)
}

//Utimens updates the access and mofication time of the file.
func (f *ParanoidFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	util.Log.Info("Utimens called on file : " + f.Name)
	var code int
	var err error
	if util.SendOverNetwork {
		if atime != nil {
			if mtime != nil {
				code, err = util.RaftServer.RequestUtimesCommand(f.Name, int64(atime.Second()), int64(atime.Nanosecond()),
					int64(mtime.Second()), int64(mtime.Nanosecond()))
			} else {
				code, err = util.RaftServer.RequestUtimesCommand(f.Name, int64(atime.Second()), int64(atime.Nanosecond()), 0, 0)
			}
		} else {
			code, err = util.RaftServer.RequestUtimesCommand(f.Name, 0, 0, int64(mtime.Second()), int64(mtime.Nanosecond()))
		}
	} else {
		code, err = commands.UtimesCommand(util.PfsDirectory, f.Name, atime, mtime)
	}

	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running utimes command :", err)
	}

	if err != nil {
		util.Log.Error("Error running utimes command :", err)
	}
	return util.GetFuseReturnCode(code)
}

//Chmod changes the permission flags of the file
func (f *ParanoidFile) Chmod(perms uint32) fuse.Status {
	util.Log.Info("Chmod called on file : " + f.Name)
	var code int
	var err error
	if util.SendOverNetwork {
		code, err = util.RaftServer.RequestChmodCommand(f.Name, perms)
	} else {
		code, err = commands.ChmodCommand(util.PfsDirectory, f.Name, os.FileMode(perms))
	}

	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running chmod command :", err)
	}

	if err != nil {
		util.Log.Error("Error running chmod command :", err)
	}
	return util.GetFuseReturnCode(code)
}
