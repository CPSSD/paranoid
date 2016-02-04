package filesystem

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/pfi/file"
	"github.com/cpssd/paranoid/pfsd/pfi/util"
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
	util.Log.Info("GetAttr called on", name)

	// Special case : "" is the root of our filesystem
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	code, err, stats := commands.StatCommand(util.PfsDirectory, name)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running stat command :", err)
	}

	if err != nil {
		util.Log.Error("Error running stat command :", err)
	}

	if code != returncodes.OK {
		return nil, util.GetFuseReturnCode(code)
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
	util.Log.Info("OpenDir called on : " + name)

	code, err, fileNames := commands.ReadDirCommand(util.PfsDirectory, name)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running readdir command :", err)
	}

	if err != nil {
		util.Log.Error("Error running readdir command :", err)
	}

	if code != returncodes.OK {
		return nil, util.GetFuseReturnCode(code)
	}

	dirEntries := make([]fuse.DirEntry, len(fileNames))
	for i, dirName := range fileNames {
		util.Log.Info("OpenDir has " + dirName)
		dirEntries[i] = fuse.DirEntry{Name: dirName}
	}

	return dirEntries, fuse.OK
}

//Open is called to get a custom file object for a certain file so that
//Read and Write (among others) opperations can be executed on this
//custom file object (ParanoidFile, see below)
func (fs *ParanoidFileSystem) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	util.Log.Info("Open called on : " + name)
	return file.NewParanoidFile(name), fuse.OK
}

//Create is called when a new file is to be created.
func (fs *ParanoidFileSystem) Create(name string, flags uint32, mode uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	util.Log.Info("Create called on : " + name)
	code, err := commands.CreatCommand(util.PfsDirectory, name, os.FileMode(mode), util.SendOverNetwork)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running creat command :", err)
	}

	if err != nil {
		util.Log.Error("Error running creat command :", err)
	}

	if code != returncodes.OK {
		return nil, util.GetFuseReturnCode(code)
	}
	return file.NewParanoidFile(name), fuse.OK
}

//Access is called by fuse to see if it has access to a certain file
func (fs *ParanoidFileSystem) Access(name string, mode uint32, context *fuse.Context) fuse.Status {
	util.Log.Info("Access called on : " + name)
	if name != "" {
		code, err := commands.AccessCommand(util.PfsDirectory, name, mode)
		if code == returncodes.EUNEXPECTED {
			util.Log.Fatal("Error running access command :", err)
		}

		if err != nil {
			util.Log.Error("Error running access command :", err)
		}
		return util.GetFuseReturnCode(code)
	}
	return fuse.OK
}

//Rename is called when renaming a file
func (fs *ParanoidFileSystem) Rename(oldName string, newName string, context *fuse.Context) fuse.Status {
	util.Log.Info("Rename called on : " + oldName + " to be renamed to " + newName)
	code, err := commands.RenameCommand(util.PfsDirectory, oldName, newName, util.SendOverNetwork)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running rename command :", err)
	}

	if err != nil {
		util.Log.Error("Error running rename command :", err)
	}
	return util.GetFuseReturnCode(code)
}

//Link creates a hard link from newName to oldName
func (fs *ParanoidFileSystem) Link(oldName string, newName string, context *fuse.Context) fuse.Status {
	util.Log.Info("Link called")
	code, err := commands.LinkCommand(util.PfsDirectory, oldName, newName, util.SendOverNetwork)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running link command :", err)
	}

	if err != nil {
		util.Log.Error("Error running link command :", err)
	}
	return util.GetFuseReturnCode(code)
}

//Symlink creates a symbolic link from newName to oldName
func (fs *ParanoidFileSystem) Symlink(oldName string, newName string, context *fuse.Context) fuse.Status {
	util.Log.Info("Symbolic link called from", oldName, "to", newName)
	code, err := commands.SymlinkCommand(util.PfsDirectory, oldName, newName, util.SendOverNetwork)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running symlink command :", err)
	}

	if err != nil {
		util.Log.Error("Error running symlink command :", err)
	}
	return util.GetFuseReturnCode(code)
}

func (fs *ParanoidFileSystem) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	util.Log.Info("Readlink called on", name)
	code, err, link := commands.ReadlinkCommand(util.PfsDirectory, name)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running readlink command :", err)
	}

	if err != nil {
		util.Log.Error("Error running readlink command :", err)
	}
	return link, util.GetFuseReturnCode(code)
}

//Unlink is called when deleting a file
func (fs *ParanoidFileSystem) Unlink(name string, context *fuse.Context) fuse.Status {
	util.Log.Info("Unlink callde on : " + name)
	code, err := commands.UnlinkCommand(util.PfsDirectory, name, util.SendOverNetwork)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running unlink command :", err)
	}

	if err != nil {
		util.Log.Error("Error running unlink command :", err)
	}
	return util.GetFuseReturnCode(code)
}

//Mkdir is called when creating a directory
func (fs *ParanoidFileSystem) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	util.Log.Info("Mkdir called on : " + name)
	code, err := commands.MkdirCommand(util.PfsDirectory, name, os.FileMode(mode), util.SendOverNetwork)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running mkdir command :", err)
	}

	if err != nil {
		util.Log.Error("Error running mkdir command :", err)
	}
	return util.GetFuseReturnCode(code)
}

//Rmdir is called when deleting a directory
func (fs *ParanoidFileSystem) Rmdir(name string, context *fuse.Context) fuse.Status {
	util.Log.Info("Rmdir called on : " + name)
	code, err := commands.RmdirCommand(util.PfsDirectory, name, util.SendOverNetwork)
	if code == returncodes.EUNEXPECTED {
		util.Log.Fatal("Error running rmdir command :", err)
	}

	if err != nil {
		util.Log.Error("Error running rmdir command :", err)
	}
	return util.GetFuseReturnCode(code)
}

//Truncate is called when a file is to be reduced in length to size.
func (fs *ParanoidFileSystem) Truncate(name string, size uint64, context *fuse.Context) fuse.Status {
	util.Log.Info("Truncate called on : " + name)
	pfile := file.NewParanoidFile(name)
	return pfile.Truncate(size)
}

//Utimens update the Access time and modified time of a given file.
func (fs *ParanoidFileSystem) Utimens(name string, atime *time.Time, mtime *time.Time, context *fuse.Context) fuse.Status {
	util.Log.Info("Utimens called on : " + name)
	pfile := file.NewParanoidFile(name)
	return pfile.Utimens(atime, mtime)
}

//Chmod is called when the permissions of a file are to be changed
func (fs *ParanoidFileSystem) Chmod(name string, perms uint32, context *fuse.Context) fuse.Status {
	util.Log.Info("Chmod called on : " + name)
	pfile := file.NewParanoidFile(name)
	return pfile.Chmod(perms)
}
