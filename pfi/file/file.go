package file

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfi/util"
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
	code, err, bytesWritten := commands.WriteCommand(util.PfsDirectory, f.Name, off, int64(len(content)), content, util.SendOverNetwork)
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
	code, err := commands.TruncateCommand(util.PfsDirectory, f.Name, int64(size), util.SendOverNetwork)
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
	code, err := commands.UtimesCommand(util.PfsDirectory, f.Name, atime, mtime, util.SendOverNetwork)
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
	code, err := commands.ChmodCommand(util.PfsDirectory, f.Name, os.FileMode(perms), util.SendOverNetwork)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running chmod command :", err)
	}

	if err != nil {
		util.Log.Error("Error running chmod command :", err)
	}
	return util.GetFuseReturnCode(code)
}
