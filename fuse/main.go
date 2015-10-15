package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

var count = 0

type HelloFs struct {
	pathfs.FileSystem
}

func (hf *HelloFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs GetAttr called")
	fmt.Println("Name : " + name)
	fmt.Println("Context : ", context)

	switch name {

	case "Mladen.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK
	case "Wojtek.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK
	case "Terry.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK
	case "Sean.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK
	case "Connor.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK
	case "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (hf *HelloFs) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs OpenDir called")
	fmt.Println("Name : " + name)
	fmt.Println("Context : ", context)

	if name == "" {
		c := []fuse.DirEntry{
			{Name: "Mladen.txt", Mode: fuse.S_IFREG},
			{Name: "Wojtek.txt", Mode: fuse.S_IFREG},
			{Name: "Terry.txt", Mode: fuse.S_IFREG},
			{Name: "Sean.txt", Mode: fuse.S_IFREG},
			{Name: "Connor.txt", Mode: fuse.S_IFREG},
		}
		return c, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (hf *HelloFs) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	count++
	fmt.Println(count, "---------------------------------")
	fmt.Println("HelloFs Open called")
	fmt.Println("Name : " + name)
	fmt.Println("Flags : ", flags)
	fmt.Println("Context : ", context)

	if name == "Mladen.txt" {
		if flags&fuse.O_ANYWRITE != 0 {
			return nil, fuse.EPERM
		}

		return nodefs.NewDataFile([]byte(name)), fuse.OK
	}
	if name == "Wojtek.txt" {
		if flags&fuse.O_ANYWRITE != 0 {
			return nil, fuse.EPERM
		}

		return nodefs.NewDataFile([]byte(name)), fuse.OK
	}
	if name == "Terry.txt" {
		if flags&fuse.O_ANYWRITE != 0 {
			return nil, fuse.EPERM
		}

		return nodefs.NewDataFile([]byte(name)), fuse.OK
	}
	if name == "Sean.txt" {
		if flags&fuse.O_ANYWRITE != 0 {
			return nil, fuse.EPERM
		}

		return nodefs.NewDataFile([]byte(name)), fuse.OK
	}
	if name == "Connor.txt" {
		if flags&fuse.O_ANYWRITE != 0 {
			return nil, fuse.EPERM
		}

		return nodefs.NewDataFile([]byte(name)), fuse.OK
	}

	return nil, fuse.ENOENT
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}
	nfs := pathfs.NewPathNodeFs(&HelloFs{FileSystem: pathfs.NewDefaultFileSystem()}, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
