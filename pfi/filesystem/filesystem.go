package filesystem

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfi/file"
	"github.com/cpssd/paranoid/pfi/pfsminterface"
	"github.com/cpssd/paranoid/pfi/util"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

//ParanoidFileSystem is the struct which holds all
//the custom filesystem functions pp2p provides
type ParanoidFileSystem struct {
	pathfs.FileSystem
}

//The structure for processing information returned by a pfs stat command
type statInfo struct {
	Length int64     `json:"length",omitempty`
	Ctime  time.Time `json:"ctime",omitempty`
	Mtime  time.Time `json:"mtime",omitempty`
	Atime  time.Time `json:"atime",omitempty`
}

//GetAttr is called by fuse when the attributes of a
//file or directory are needed. (pfs stat)
func (fs *ParanoidFileSystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	util.LogMessage("GetAttr called on : " + name)

	// Special case : "" is the root of our flat
	// file system (Only directory GetAttr can be called on)
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	code, output := pfsminterface.RunCommand(nil, "stat", util.PfsDirectory, name)
	if code == pfsminterface.ENOENT {
		util.LogMessage("stat file doesn't exist " + name)
		return nil, fuse.ENOENT
	}

	stats := statInfo{}
	err := json.Unmarshal(output, &stats)
	if err != nil {
		log.Fatalln("Error processing JSON returned by stat command:", err)
	}

	attr := fuse.Attr{
		Size:  uint64(stats.Length),
		Atime: uint64(stats.Atime.Unix()),
		Ctime: uint64(stats.Ctime.Unix()),
		Mtime: uint64(stats.Mtime.Unix()),
		Mode:  fuse.S_IFREG | 0644, // S_IFREG = regular file
	}

	return &attr, fuse.OK
}

//OpenDir is called when the contents of a directory are needed. There
//is only one directory this can be called on in this sprint. i.e
//the root directory of the file system.
func (fs *ParanoidFileSystem) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	util.LogMessage("OpenDir called on : " + name)
	code, output := pfsminterface.RunCommand(nil, "readdir", util.PfsDirectory)
	if code == pfsminterface.ENOENT {
		return nil, fuse.ENOENT
	}
	outputString := string(output)

	util.LogMessage("OpenDir returns : " + outputString)
	if outputString == "" {
		dirEntries := make([]fuse.DirEntry, 0)
		return dirEntries, fuse.OK
	}

	fileNames := strings.Split(outputString, "\n")
	dirEntries := make([]fuse.DirEntry, len(fileNames))

	for i, dirName := range fileNames {
		util.LogMessage("OpenDir has " + dirName)
		dirEntries[i] = fuse.DirEntry{Name: dirName}
	}

	return dirEntries, fuse.OK
}

//Open is called to get a custom file object for a certain file so that
//Read and Write (among others) opperations can be executed on this
//custom file object (ParanoidFile, see below)
func (fs *ParanoidFileSystem) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	util.LogMessage("Open called on : " + name)
	return file.NewParanoidFile(name), fuse.OK
}

//Create is called when a new file is to be created.
func (fs *ParanoidFileSystem) Create(name string, flags uint32, mode uint32, context *fuse.Context) (retfile nodefs.File, code fuse.Status) {
	util.LogMessage("Create called on : " + name)
	retcode, _ := pfsminterface.RunCommand(nil, "creat", util.PfsDirectory, name, strconv.Itoa(int(mode)))
	if retcode == pfsminterface.ENOENT {
		return nil, fuse.ENOENT
	}
	retfile = file.NewParanoidFile(name)
	return retfile, fuse.OK
}

//Access is called by fuse to see if it has access to a certain file
func (fs *ParanoidFileSystem) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Access called on : " + name)
	if name != "" {
		retcode, _ := pfsminterface.RunCommand(nil, "access", util.PfsDirectory, name, strconv.Itoa(int(mode)))
		if retcode == pfsminterface.ENOENT {
			return fuse.ENOENT
		} else if retcode == pfsminterface.EACCES {
			return fuse.EACCES
		}
	}
	return fuse.OK
}

//Rename is called when renaming a file
func (fs *ParanoidFileSystem) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Rename called on : " + oldName + " to be renamed to " + newName)
	retcode, _ := pfsminterface.RunCommand(nil, "rename", util.PfsDirectory, oldName, newName)

	if retcode == pfsminterface.ENOENT {
		return fuse.ENOENT
	} else if retcode == pfsminterface.EACCES {
		return fuse.EACCES
	}

	return fuse.OK
}

//Write content to the given file starting at off
func (fs *ParanoidFileSystem) Write(name string, content []byte, off int64, context *fuse.Context) (uint32, fuse.Status) {
	util.LogMessage("Write called on : " + name)
	pfile := file.NewParanoidFile(name)
	return pfile.Write(content, off)
}

//Read a given file starting at off.
func (fs *ParanoidFileSystem) Read(name string, buf []byte, off int64, context *fuse.Context) (fuse.ReadResult, fuse.Status) {
	util.LogMessage("Read called on : " + name)
	pfile := file.NewParanoidFile(name)
	return pfile.Read(buf, off)
}

//Truncate is called when a file is to be reduced in length to size.
func (fs *ParanoidFileSystem) Truncate(name string, size uint64, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Truncate called on : " + name)
	pfile := file.NewParanoidFile(name)
	return pfile.Truncate(size)
}

//Utimens update the Acess time and modified time of a given file.
func (fs *ParanoidFileSystem) Utimens(name string, atime *time.Time, mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Utimens called on : " + name)
	pfile := file.NewParanoidFile(name)
	return pfile.Utimens(atime, mtime)
}

//Chmod is called when the permissions of a file are to be changed
func (fs *ParanoidFileSystem) Chmod(name string, perms uint32, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Chmod called on : " + name)
	pfile := file.NewParanoidFile(name)
	return pfile.Chmod(perms)
}
