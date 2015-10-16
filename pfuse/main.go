package main

import (
	"fmt"
	"github.com/cpssd/paranoid/pfuse/pfs_interface"
	"log"
	"os"

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

	stats := pfsInterface.Stat(pfsInitPoint, pfsLocation, name)
	attr := fuse.Attr{
		Size:  uint64(stats.Length),
		Atime: uint64(stats.Atime.Unix()),
		Ctime: uint64(stats.Ctime.Unix()),
		Mtime: uint64(stats.Mtime.Unix()),
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
	fmt.Println("Here")
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

	data := pfsInterface.Read(pfsInitPoint, pfsLocation, name, -1, -1)
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	return nodefs.NewDataFile(data), fuse.OK
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
	file = nodefs.NewDefaultFile()
	return file, fuse.OK
}
