package pfi

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pfi/glob"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"os"
	"time"
)

//ParanoidFileSystem is the struct which holds all
//the custom filesystem functions pp2p provides
type ParanoidFileSystem struct {
	pathfs.FileSystem
}

//GetAttr is called by fuse when the attributes of a
//file or directory are needed. (pfs stat)
func (fs *ParanoidFileSystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	Log.Verbose("GetAttr called on", name)

	// Special case : "" is the root of our filesystem
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	code, err, stats := commands.StatCommand(globals.ParanoidDir, name)
	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running stat command :", err)
	}

	if err != nil { // TODO this produces tonnes of usless logspam
		Log.Error("Error running stat command :", err)
	}

	if code != returncodes.OK {
		return nil, GetFuseReturnCode(code)
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
	Log.Verbose("OpenDir called on : " + name)

	code, err, fileNames := commands.ReadDirCommand(globals.ParanoidDir, name)
	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running readdir command :", err)
	}

	if err != nil {
		Log.Error("Error running readdir command :", err)
	}

	if code != returncodes.OK {
		return nil, GetFuseReturnCode(code)
	}

	dirEntries := make([]fuse.DirEntry, len(fileNames))
	for i, dirName := range fileNames {
		Log.Verbose("OpenDir has " + dirName)
		dirEntries[i] = fuse.DirEntry{Name: dirName}
	}

	return dirEntries, fuse.OK
}

//Open is called to get a custom file object for a certain file so that
//Read and Write (among others) opperations can be executed on this
//custom file object (ParanoidFile, see below)
func (fs *ParanoidFileSystem) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	Log.Verbose("Open called on : " + name)
	return newParanoidFile(name), fuse.OK
}

//Create is called when a new file is to be created.
func (fs *ParanoidFileSystem) Create(name string, flags uint32, mode uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	Log.Info("Create called on : " + name)
	var code returncodes.Code
	var err error
	shouldGlob := glob.ShouldIgnore(name, false)
	if SendOverNetwork && !shouldGlob {
		code, err = globals.RaftNetworkServer.RequestCreatCommand(name, mode)
	} else {
		code, err = commands.CreatCommand(globals.ParanoidDir, name, os.FileMode(mode), shouldGlob)
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running creat command :", err)
	}

	if err != nil {
		Log.Error("Error running creat command :", err)
	}

	if code != returncodes.OK {
		return nil, GetFuseReturnCode(code)
	}
	return newParanoidFile(name), fuse.OK
}

//Access is called by fuse to see if it has access to a certain file
func (fs *ParanoidFileSystem) Access(name string, mode uint32, context *fuse.Context) fuse.Status {
	Log.Verbose("Access called on : " + name)
	if name != "" {
		code, err := commands.AccessCommand(globals.ParanoidDir, name, mode)
		if code == returncodes.EUNEXPECTED {
			Log.Fatal("Error running access command :", err)
		}

		if err != nil {
			Log.Error("Error running access command :", err)
		}
		return GetFuseReturnCode(code)
	}
	return fuse.OK
}

//Rename is called when renaming a file
func (fs *ParanoidFileSystem) Rename(oldName string, newName string, context *fuse.Context) fuse.Status {
	Log.Info("Rename called on : " + oldName + " to be renamed to " + newName)
	var code returncodes.Code
	var err error
	if SendOverNetwork && !glob.ShouldIgnore(newName, false) {
		code, err = globals.RaftNetworkServer.RequestRenameCommand(oldName, newName)
	} else {
		code, err = commands.RenameCommand(globals.ParanoidDir, oldName, newName)
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running rename command :", err)
	}

	if err != nil {
		Log.Error("Error running rename command :", err)
	}
	return GetFuseReturnCode(code)
}

//Link creates a hard link from newName to oldName
func (fs *ParanoidFileSystem) Link(oldName string, newName string, context *fuse.Context) fuse.Status {
	Log.Info("Link called")
	var code returncodes.Code
	var err error
	if SendOverNetwork && !glob.ShouldIgnore(newName, false) {
		code, err = globals.RaftNetworkServer.RequestLinkCommand(oldName, newName)
	} else {
		code, err = commands.LinkCommand(globals.ParanoidDir, oldName, newName)
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running link command :", err)
	}

	if err != nil {
		Log.Error("Error running link command :", err)
	}
	return GetFuseReturnCode(code)
}

//Symlink creates a symbolic link from newName to oldName
func (fs *ParanoidFileSystem) Symlink(oldName string, newName string, context *fuse.Context) fuse.Status {
	Log.Info("Symbolic link called from", oldName, "to", newName)
	var code returncodes.Code
	var err error
	if SendOverNetwork && !glob.ShouldIgnore(newName, false) {
		code, err = globals.RaftNetworkServer.RequestSymlinkCommand(oldName, newName)
	} else {
		code, err = commands.SymlinkCommand(globals.ParanoidDir, oldName, newName)
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running symlink command :", err)
	}

	if err != nil {
		Log.Error("Error running symlink command :", err)
	}
	return GetFuseReturnCode(code)
}

func (fs *ParanoidFileSystem) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	Log.Info("Readlink called on", name)
	code, err, link := commands.ReadlinkCommand(globals.ParanoidDir, name)
	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running readlink command :", err)
	}

	if err != nil {
		Log.Error("Error running readlink command :", err)
	}
	return link, GetFuseReturnCode(code)
}

//Unlink is called when deleting a file
func (fs *ParanoidFileSystem) Unlink(name string, context *fuse.Context) fuse.Status {
	Log.Info("Unlink called on : " + name)
	var code returncodes.Code
	var err error
	if SendOverNetwork && !glob.ShouldIgnore(name, false) {
		code, err = globals.RaftNetworkServer.RequestUnlinkCommand(name)
	} else {
		code, err = commands.UnlinkCommand(globals.ParanoidDir, name)
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running unlink command :", err)
	}

	if err != nil {
		Log.Error("Error running unlink command :", err)
	}
	return GetFuseReturnCode(code)
}

//Mkdir is called when creating a directory
func (fs *ParanoidFileSystem) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	Log.Info("Mkdir called on : " + name)
	var code returncodes.Code
	var err error
	if SendOverNetwork && !glob.ShouldIgnore(name, false) {
		code, err = globals.RaftNetworkServer.RequestMkdirCommand(name, mode)
	} else {
		code, err = commands.MkdirCommand(globals.ParanoidDir, name, os.FileMode(mode))
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running mkdir command :", err)
	}

	if err != nil {
		Log.Error("Error running mkdir command :", err)
	}
	return GetFuseReturnCode(code)
}

//Rmdir is called when deleting a directory
func (fs *ParanoidFileSystem) Rmdir(name string, context *fuse.Context) fuse.Status {
	Log.Info("Rmdir called on : " + name)
	var code returncodes.Code
	var err error
	if SendOverNetwork && !glob.ShouldIgnore(name, true) {
		code, err = globals.RaftNetworkServer.RequestRmdirCommand(name)
	} else {
		code, err = commands.RmdirCommand(globals.ParanoidDir, name)
	}

	if code == returncodes.EUNEXPECTED {
		Log.Fatal("Error running rmdir command :", err)
	}

	if err != nil {
		Log.Error("Error running rmdir command :", err)
	}
	return GetFuseReturnCode(code)
}

//Truncate is called when a file is to be reduced in length to size.
func (fs *ParanoidFileSystem) Truncate(name string, size uint64, context *fuse.Context) fuse.Status {
	Log.Info("Truncate called on : " + name)
	pfile := newParanoidFile(name)
	return pfile.Truncate(size)
}

//Utimens update the Access time and modified time of a given file.
func (fs *ParanoidFileSystem) Utimens(name string, atime *time.Time, mtime *time.Time, context *fuse.Context) fuse.Status {
	Log.Info("Utimens called on : " + name)
	pfile := newParanoidFile(name)
	return pfile.Utimens(atime, mtime)
}

//Chmod is called when the permissions of a file are to be changed
func (fs *ParanoidFileSystem) Chmod(name string, perms uint32, context *fuse.Context) fuse.Status {
	Log.Info("Chmod called on : " + name)
	pfile := newParanoidFile(name)
	return pfile.Chmod(perms)
}
