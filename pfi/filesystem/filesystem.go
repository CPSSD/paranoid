package filesystem

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfi/file"
	"github.com/cpssd/paranoid/pfi/pfsminterface"
	"github.com/cpssd/paranoid/pfi/util"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//ParanoidFileSystem is the struct which holds all
//the custom filesystem functions pp2p provides
type ParanoidFileSystem struct {
	pathfs.FileSystem
}

//The structure for processing information returned by a pfs stat command
type statInfo struct {
	Length int64       `json:"length",omitempty`
	Ctime  time.Time   `json:"ctime",omitempty`
	Mtime  time.Time   `json:"mtime",omitempty`
	Atime  time.Time   `json:"atime",omitempty`
	Mode   os.FileMode `json:"mode",omitempty`
}

//GetAttr is called by fuse when the attributes of a
//file or directory are needed. (pfs stat)
func (fs *ParanoidFileSystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	util.LogMessage("GetAttr called on : " + name)

	// Special case : "" is the root of our filesystem
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	code, output := pfsminterface.RunCommand(nil, "stat", util.PfsDirectory, name)
	if code != returncodes.OK {
		return nil, util.GetFuseReturnCode(code)
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
		Mode:  uint32(stats.Mode),
	}

	return &attr, fuse.OK
}

//OpenDir is called when the contents of a directory are needed.
func (fs *ParanoidFileSystem) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	util.LogMessage("OpenDir called on : " + name)

	code, output := pfsminterface.RunCommand(nil, "readdir", util.PfsDirectory, name)
	if code != returncodes.OK {
		return nil, util.GetFuseReturnCode(code)
	}
	outputString := string(output)

	util.LogMessage("OpenDir returns : " + outputString)
	if outputString == "" {
		var dirEntries []fuse.DirEntry
		return dirEntries, fuse.OK
	}

	fileNames := strings.Split(outputString, "\n")
	fileNames = fileNames[:len(fileNames)-1]
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
	retcode, _ := pfsminterface.RunCommand(nil, "creat", util.PfsDirectory, name, strconv.FormatInt(int64(mode), 8))
	if retcode != returncodes.OK {
		return nil, util.GetFuseReturnCode(retcode)
	}
	retfile = file.NewParanoidFile(name)
	return retfile, fuse.OK
}

//Access is called by fuse to see if it has access to a certain file
func (fs *ParanoidFileSystem) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Access called on : " + name)
	if name != "" {
		retcode, _ := pfsminterface.RunCommand(nil, "access", util.PfsDirectory, name, strconv.Itoa(int(mode)))
		return util.GetFuseReturnCode(retcode)
	}
	return fuse.OK
}

//Rename is called when renaming a file
func (fs *ParanoidFileSystem) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Rename called on : " + oldName + " to be renamed to " + newName)
	retcode, _ := pfsminterface.RunCommand(nil, "rename", util.PfsDirectory, oldName, newName)
	return util.GetFuseReturnCode(retcode)
}

//Link creates a hard link from newName to oldName
func (fs *ParanoidFileSystem) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Link called")
	retcode, _ := pfsminterface.RunCommand(nil, "link", util.PfsDirectory, oldName, newName)
	return util.GetFuseReturnCode(retcode)
}

//Unlink is called when deleting a file
func (fs *ParanoidFileSystem) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Unlink callde on : " + name)
	retcode, _ := pfsminterface.RunCommand(nil, "unlink", util.PfsDirectory, name)
	return util.GetFuseReturnCode(retcode)
}

//Mkdir is called when creating a directory
func (fs *ParanoidFileSystem) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	util.LogMessage("Mkdir called on : " + name)
	retcode, _ := pfsminterface.RunCommand(nil, "mkdir", util.PfsDirectory, name, strconv.FormatInt(int64(mode), 8))
	return util.GetFuseReturnCode(retcode)
}

//Rmdir is called when deleting a directory
func (fs *ParanoidFileSystem) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	util.LogMessage("Rmdir called on : " + name)
	retcode, _ := pfsminterface.RunCommand(nil, "rmdir", util.PfsDirectory, name)
	return util.GetFuseReturnCode(retcode)
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
