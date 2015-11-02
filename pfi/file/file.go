package file

import (
	"github.com/cpssd/paranoid/pfi/pfsminterface"
	"github.com/cpssd/paranoid/pfi/util"
	"strconv"

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
	util.LogMessage("Read called on : " + f.Name)
	code, data := pfsminterface.RunCommand(nil, "read", util.PfsDirectory, f.Name, strconv.FormatInt(off, 10), strconv.FormatInt(int64(len(buf)), 10))
	if code == pfsminterface.ENOENT {
		return nil, fuse.ENOENT
	}
	return fuse.ReadResultData(data), fuse.OK
}

//Write writes to a file
func (f *ParanoidFile) Write(content []byte, off int64) (uint32, fuse.Status) {
	util.LogMessage("Write called on : " + f.Name)
	code, _ := pfsminterface.RunCommand(content, "write", util.PfsDirectory, f.Name, strconv.FormatInt(off, 10), strconv.FormatInt(int64(len(content)), 10))
	if code == pfsminterface.ENOENT {
		return 0, fuse.ENOENT
	}
	return uint32(len(content)), fuse.OK
}

//Truncate is called when a file is to be reduced in length to size.
func (f *ParanoidFile) Truncate(size uint64) fuse.Status {
	util.LogMessage("Truncate called on file : " + f.Name)
	code, _ := pfsminterface.RunCommand(nil, "truncate", util.PfsDirectory, f.Name, strconv.FormatInt(int64(size), 10))
	if code == pfsminterface.ENOENT {
		return fuse.ENOENT
	}
	return fuse.OK
}

//Chmod changes the permission flags of the file
func (f *ParanoidFile) Chmod(perms uint32) fuse.Status {
	util.LogMessage("Chmod called on file : " + f.Name)
	code, _ := pfsminterface.RunCommand(nil, "chmod", util.PfsDirectory, f.Name, strconv.FormatInt(int64(perms), 10))
	if code == pfsminterface.ENOENT {
		return fuse.ENOENT
	}
	return fuse.OK
}
