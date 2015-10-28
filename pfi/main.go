package main

import (
	"encoding/json"
	"flag"
	"github.com/cpssd/paranoid/pfi/pfsinterface"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

var mountPoint string
var pfsInitPoint string
var logOutput *bool

func main() {
	// parsing flags and args
	logOutput = flag.Bool("v", false, "Log operations in standard output")
	markNetwork := flag.Bool("n", false, "Mark file system operations as coming from the network")
	flag.Parse()
	if *markNetwork {
		pfsinterface.OriginFlag = "-n"
	} else {
		pfsinterface.OriginFlag = "-f"
	}
	noFlagArgs := flag.Args()

	if len(noFlagArgs) < 2 {
		log.Fatalln("\nUsage:\npfi [flags] <PfsInitPoint> <MountPoint>")
	}

	var err error
	pfsInitPoint, err = filepath.Abs(noFlagArgs[0])
	if err != nil {
		log.Fatalln(err)
	}
	mountPoint, err = filepath.Abs(noFlagArgs[1])
	if err != nil {
		log.Fatalln(err)
	}

	// configuring log
	log.SetFlags(log.Ldate | log.Ltime)

	// setting up with fuse
	nfs := pathfs.NewPathNodeFs(&ParanoidFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
	}, nil)
	server, _, err := nodefs.MountRoot(mountPoint, nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}

//ParanoidFileSystem is the struct which holds all
//the custom filesystem functions pp2p provides
type ParanoidFileSystem struct {
	pathfs.FileSystem
}

//The structure for processing information returned by a pfs stat command
type statInfo struct {
	Exists bool      `json:"exists",omitempty`
	Length int64     `json:"length",omitempty`
	Ctime  time.Time `json:"ctime",omitempty`
	Mtime  time.Time `json:"mtime",omitempty`
	Atime  time.Time `json:"atime",omitempty`
}

//GetAttr is called by fuse when the attributes of a
//file or directory are needed. (pfs stat)
func (fs *ParanoidFileSystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	logMessage("GetAttr called on : " + name)

	// Special case : "" is the root of our flat
	// file system (Only directory GetAttr can be called on)
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	output := pfsinterface.RunCommand(nil, "stat", pfsInitPoint, name)
	stats := statInfo{}
	err := json.Unmarshal(output, &stats)
	if err != nil {
		log.Fatalln("Error processing JSON returned by stat command:", err)
	}

	if stats.Exists == false {
		logMessage("stat file doesn't exist " + name)
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

//OpenDir is called when the contents of a directory are needed. There
//is only one directory this can be called on in this sprint. i.e
//the root directory of the file system.
func (fs *ParanoidFileSystem) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	logMessage("OpenDir called on : " + name)
	// pfs init point is used instead of a name in pfsinterface.Readdir indicating
	// the root of the file system (because name == "")
	output := pfsinterface.RunCommand(nil, "readdir", pfsInitPoint)
	outputString := string(output)

	logMessage("OpenDir returns : " + outputString)
	if outputString == "" {
		dirEntries := make([]fuse.DirEntry, 0)
		return dirEntries, fuse.OK
	}

	fileNames := strings.Split(outputString, "\n")
	dirEntries := make([]fuse.DirEntry, len(fileNames))

	for i, dirName := range fileNames {
		logMessage("OpenDir has " + dirName)
		dirEntries[i] = fuse.DirEntry{Name: dirName}
	}

	return dirEntries, fuse.OK
}

//Open is called to get a custom file object for a certain file so that
//Read and Write (among others) opperations can be executed on this
//custom file object (ParanoidFile, see below)
func (fs *ParanoidFileSystem) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	logMessage("Open called on : " + name)
	return NewParanoidFile(name), fuse.OK
}

//Create is called when a new file is to be created.
func (fs *ParanoidFileSystem) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	logMessage("Create called on : " + name)
	pfsinterface.RunCommand(nil, "creat", pfsInitPoint, name)
	file = NewParanoidFile(name)
	return file, fuse.OK
}

//Access is called by fuse to see if it has access to a certain
//file. In this sprint access will always be granted.
func (fs *ParanoidFileSystem) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	logMessage("Access called on : " + name)
	return fuse.OK
}

//Truncate is called when a file is to be truncated. We dont have this functionality
//yet but its added in here so that if fuse calls it it doesn't break the program.
func (fs *ParanoidFileSystem) Truncate(name string, size uint64, context *fuse.Context) (code fuse.Status) {
	logMessage("Truncate called on : " + name)
	return fuse.OK
}

//Utimens : We dont have this functionality implemented but
//sometimes when calling callind touch on some file
//fuse leaves behind an annoying message :
//"touch: setting times of ‘filename’: Function not implemented"
//I added this function that just returns OK to suppress the annoying message.
func (fs *ParanoidFileSystem) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	logMessage("Utimens called on : " + name)
	return fuse.OK
}

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
	logMessage("Read called on : " + f.Name)
	data := pfsinterface.RunCommand(nil, "read", pfsInitPoint, f.Name, strconv.FormatInt(off, 10), strconv.FormatInt(int64(len(buf)), 10))
	return fuse.ReadResultData(data), fuse.OK
}

//Write writes to a file
func (f *ParanoidFile) Write(content []byte, off int64) (uint32, fuse.Status) {
	logMessage("Write called on : " + f.Name)
	pfsinterface.RunCommand(content, "write", pfsInitPoint, strconv.FormatInt(off, 10), strconv.FormatInt(int64(len(content)), 10))
	return uint32(len(content)), fuse.OK
}

//LogMessage checks if the -v flag was specified and either logs or doesnt log the message
func logMessage(message string) {
	if *logOutput {
		log.Println(message)
	}
}
