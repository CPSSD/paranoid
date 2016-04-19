package pfi

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pfi/glob"
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

//ParanoidFile is a custom file struct with read and write functions
type ParanoidFile struct {
	Name string
	nodefs.File
}

//newParanoidFile returns a new object of ParanoidFile
func newParanoidFile(name string) nodefs.File {
	return &ParanoidFile{
		Name: name,
		File: nodefs.NewDefaultFile(),
	}
}

//Read reads a file and returns an array of bytes
func (f *ParanoidFile) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {
	Log.Verbose("Read called on file:", f.Name)

	code, err, data := commands.ReadCommand(globals.ParanoidDir, f.Name, off, int64(len(buf)))
	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running read command :", err)
	}

	if err != nil {
		Log.Error("Error running read command :", err)
	}

	copy(buf, data)
	if code != returncodes.OK {
		return nil, GetFuseReturnCode(code)
	}
	return fuse.ReadResultData(data), fuse.OK
}

//Write writes to a file
func (f *ParanoidFile) Write(content []byte, off int64) (uint32, fuse.Status) {
	Log.Info("Write called on file : " + f.Name)
	var (
		code         returncodes.Code
		err          error
		bytesWritten int
	)
	changeInode := false
	if SendOverNetwork && !glob.ShouldIgnore(f.Name, changeInode) {
		code, err, bytesWritten = globals.RaftNetworkServer.RequestWriteCommand(f.Name, off, int64(len(content)), content)
	} else {
		code, err, bytesWritten = commands.WriteCommand(globals.ParanoidDir, f.Name, off, int64(len(content)), content)
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running write command :", err)
	}

	if err != nil {
		Log.Error("Error running write command :", err)
	}

	if code != returncodes.OK {
		return 0, GetFuseReturnCode(code)
	}

	return uint32(bytesWritten), fuse.OK
}

//Truncate is called when a file is to be reduced in length to size.
func (f *ParanoidFile) Truncate(size uint64) fuse.Status {
	Log.Info("Truncate called on file : " + f.Name)
	Log.Info("TRUCATE SIZE:", size)
	var code returncodes.Code
	var err error
	changeInode := false
	if size <= 0 {
		changeInode = true
	}
	if SendOverNetwork && !glob.ShouldIgnore(f.Name, changeInode) {
		code, err = globals.RaftNetworkServer.RequestTruncateCommand(f.Name, int64(size))
	} else {
		code, err = commands.TruncateCommand(globals.ParanoidDir, f.Name, int64(size))
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running truncate command :", err)
	}

	if err != nil {
		Log.Error("Error running truncate command :", err)
	}

	return GetFuseReturnCode(code)
}

//Utimens updates the access and mofication time of the file.
func (f *ParanoidFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	Log.Info("Utimens called on file : " + f.Name)
	var code returncodes.Code
	var err error
	if SendOverNetwork && !glob.ShouldIgnore(f.Name, false) {
		code, err = globals.RaftNetworkServer.RequestUtimesCommand(f.Name, atime, mtime)
	} else {
		code, err = commands.UtimesCommand(globals.ParanoidDir, f.Name, atime, mtime)
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running utimes command :", err)
	}

	if err != nil {
		Log.Error("Error running utimes command :", err)
	}
	return GetFuseReturnCode(code)
}

//Chmod changes the permission flags of the file
func (f *ParanoidFile) Chmod(perms uint32) fuse.Status {
	Log.Info("Chmod called on file : " + f.Name)
	var code returncodes.Code
	var err error
	if SendOverNetwork && !glob.ShouldIgnore(f.Name, false) {
		code, err = globals.RaftNetworkServer.RequestChmodCommand(f.Name, perms)
	} else {
		code, err = commands.ChmodCommand(globals.ParanoidDir, f.Name, os.FileMode(perms))
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running chmod command :", err)
	}

	if err != nil {
		Log.Error("Error running chmod command :", err)
	}
	return GetFuseReturnCode(code)
}
