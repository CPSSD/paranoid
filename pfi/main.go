package main

import (
	"fmt"
	"github.com/cpssd/paranoid/pfi/pfs_interface"
	"log"
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

var count int
var mountPoint string
var pfsLocation string
var pfsInitPoint string

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		log.Fatal("\nUsage:\npfuse Mountpoint PfsInitPoint PfsExeLocation")
	}
	mountPoint, pfsInitPoint, pfsLocation = args[0], args[1], args[2]
	nfs := pathfs.NewPathNodeFs(&HelloFs{FileSystem: pathfs.NewDefaultFileSystem()}, nil)
	server, _, err := nodefs.MountRoot(mountPoint, nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}

// filesystem operations
type HelloFs struct {
	pathfs.FileSystem
}

func (hf *HelloFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs GetAttr called")
	fmt.Println("Name : " + name)
	fmt.Println("Context : ", context)

	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	stats, err := pfsInterface.Stat(pfsInitPoint, pfsLocation, name)

	if err != nil {
		return nil, fuse.ENOENT
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

func (hf *HelloFs) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs OpenDir called")
	fmt.Println("Name : " + name)
	fmt.Println("Context : ", context)

	fileNames := pfsInterface.Readdir(pfsInitPoint, pfsLocation, pfsInitPoint) // TODO: Change this back to name
	dirEntries := make([]fuse.DirEntry, len(fileNames))

	for i, dirName := range fileNames {
		dirEntries[i] = fuse.DirEntry{Name: dirName}
	}

	return dirEntries, fuse.OK
}

func (hf *HelloFs) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs Open called")
	fmt.Println("Name : " + name)
	fmt.Println("Flags : ", flags)
	fmt.Println("Context : ", context)
	file := NewPfile(name)
	return file, fuse.OK
}

func (hf *HelloFs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs Create called")
	fmt.Println("Name : " + name)
	fmt.Println("Flags : ", flags)
	fmt.Println("Mode : ", mode)
	fmt.Println("Context : ", context)

	pfsInterface.Creat(pfsInitPoint, pfsLocation, name)
	file = NewPfile(name)
	return file, fuse.OK
}

func (hf *HelloFs) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs Access called")
	fmt.Println("Name : " + name)
	fmt.Println("Mode : ", mode)
	fmt.Println("Context : ", context)
	return fuse.OK
}

func (hf *HelloFs) Truncate(name string, size uint64, context *fuse.Context) (code fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs Truncate called")
	fmt.Println("Name : " + name)
	fmt.Println("Size : ", size)
	fmt.Println("Context : ", context)
	return fuse.OK
}

func (hf *HelloFs) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs Utimes called")
	fmt.Println("Name : ", name)
	fmt.Println("Atime : ", Atime)
	fmt.Println("Mtime : ", Mtime)
	fmt.Println("Context : ", context)

	/*
		// temporary fix while file structure is flat
		filenames := pfsInterface.Readdir(pfsInitPoint, pfsLocation, pfsInitPoint)
		for _, fname := range filenames {
			if fname == name {
				// file exists
				// write nothing so times will be updated
				pfsInterface.Write(pfsInitPoint, pfsLocation, name, make([]byte, 0), 0, 0)
				return fuse.OK
			}
		}
		// file doesn't exists
		// create file
		pfsInterface.Creat(pfsInitPoint, pfsLocation, name)*/
	pfsInterface.Write(pfsInitPoint, pfsLocation, name, make([]byte, 0), 0, 0)
	return fuse.OK
}

// file opperations
type Pfile struct {
	Name string
	nodefs.File
}

func NewPfile(name string) nodefs.File {
	return &Pfile{
		Name: name,
		File: nodefs.NewDefaultFile(),
	}
}

func (f *Pfile) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("pfile Read Called on ", f.Name)

	data := pfsInterface.Read(pfsInitPoint, pfsLocation, f.Name, off, int64(len(buf)))
	return fuse.ReadResultData(data), fuse.OK
}

func (f *Pfile) Write(content []byte, off int64) (uint32, fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	pfsInterface.Write(pfsInitPoint, pfsLocation, f.Name, content, off, int64(len(content)))
	return uint32(len(content)), fuse.OK
}
func (f *Pfile) Utimens(a *time.Time, m *time.Time) fuse.Status {
	return fuse.OK
}
