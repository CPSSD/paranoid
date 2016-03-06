// +build integration

package pfi

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"os"
	"path"
	"testing"
	"time"
)

func createTestDir(t *testing.T, name string) {
	os.RemoveAll(path.Join(os.TempDir(), name))
	err := os.Mkdir(path.Join(os.TempDir(), name), 0777)
	if err != nil {
		t.Error(err)
	}
}

func removeTestDir(name string) {
	os.RemoveAll(path.Join(os.TempDir(), name))
}

func TestMain(m *testing.M) {
	Log = logger.New("testPackage", "testComponent", os.DevNull)
	globals.ParanoidDir = path.Join(os.TempDir(), "pfiTestPfsDir")
	commands.Log = logger.New("testPackage", "testComponent", os.DevNull)
	os.Exit(m.Run())
}

func setuptesting(t *testing.T) {
	removeTestDir("pfiTestPfsDir")
	createTestDir(t, "pfiTestPfsDir")
	_, err := commands.InitCommand(path.Join(os.TempDir(), "pfiTestPfsDir"))
	if err != nil {
		Log.Fatal("Error initing paranoid file system:", err)
	}
}

func TestFuseFilePerms(t *testing.T) {
	setuptesting(t)
	defer removeTestDir("pfiTestPfsDir")

	_, err := commands.CreatCommand(path.Join(os.TempDir(), "pfiTestPfsDir"), "helloworld.txt", os.FileMode(0777))
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}

	pfs := &ParanoidFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
	}

	attr, status := pfs.GetAttr("helloworld.txt", nil)
	if status != fuse.OK {
		t.Error("Error calling GetAttr")
	}
	if os.FileMode(attr.Mode).Perm() != 0777 {
		t.Error("Recieved incorrect permisions", os.FileMode(attr.Mode))
	}

	canAccess := pfs.Access("helloworld.txt", 4, nil)
	if canAccess != fuse.OK {
		t.Error("Should be able to access file")
	}

	code := pfs.Chmod("helloworld.txt", uint32(os.FileMode(0377)), nil)
	if code != fuse.OK {
		t.Error("Chmod failed error : ", code)
	}

	canAccess = pfs.Access("helloworld.txt", 4, nil)
	if canAccess != fuse.EACCES {
		t.Error("Should not be able to access file error :", canAccess)
	}
}

func TestFuseFileOperations(t *testing.T) {
	setuptesting(t)
	defer removeTestDir("pfiTestPfsDir")

	pfs := &ParanoidFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
	}

	file, code := pfs.Create("helloworld.txt", 0, uint32(os.FileMode(0777)), nil)
	if code != fuse.OK {
		t.Error("Failed to create file error : ", code)
	}

	_, code = file.Write([]byte("TEST"), 0)
	if code != fuse.OK {
		t.Error("Failed to write to file error : ", code)
	}

	buf := make([]byte, 4)
	readRes, code := file.Read(buf, 0)
	if code != fuse.OK {
		t.Error("Failed to read file error : ", code)
	}
	data, code := readRes.Bytes(buf)
	if code != fuse.OK {
		t.Error("Failed to read file error : ", code)
	}

	if string(data) != "TEST" {
		t.Error("Data read from file is not correct. Actual : ", data)
	}

	code = file.Truncate(2)
	if code != fuse.OK {
		t.Error("Failed to truncate file error : ", code)
	}

	buf = make([]byte, 2)
	readRes, code = file.Read(buf, 0)
	if code != fuse.OK {
		t.Error("Failed to read file error : ", code)
	}
	data, code = readRes.Bytes(buf)
	if code != fuse.OK {
		t.Error("Failed to read file error : ", code)
	}

	if string(data) != "TE" {
		t.Error("Data read from file is not correct. Actual : ", data)
	}
}

func TestFuseFileSystemOperations(t *testing.T) {
	setuptesting(t)
	defer removeTestDir("pfiTestPfsDir")

	pfs := &ParanoidFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
	}

	_, code := pfs.Create("file1", 0, uint32(os.FileMode(0777)), nil)
	if code != fuse.OK {
		t.Error("Failed to create new file, error : ", code)
	}

	dirEntries, code := pfs.OpenDir("", nil)
	if code != fuse.OK {
		t.Error("Could not open directory, error : ", code)
	}
	if len(dirEntries) != 1 {
		t.Error("Incorrect number of files in directory : ", dirEntries)
	}
	if dirEntries[0].Name != "file1" {
		t.Error("Incorrect file name recieved : ", dirEntries[0].Name)
	}

	code = pfs.Rename("file1", "file2", nil)
	if code != fuse.OK {
		t.Error("Failed to rename file, error : ", code)
	}

	dirEntries, code = pfs.OpenDir("", nil)
	if code != fuse.OK {
		t.Error("Could not open directory, error : ", code)
	}
	if len(dirEntries) != 1 {
		t.Error("Incorrect number of files in directory : ", len(dirEntries))
	}
	if dirEntries[0].Name != "file2" {
		t.Error("Incorrect file name recieved : ", dirEntries[0].Name)
	}

	code = pfs.Unlink("file2", nil)
	if code != fuse.OK {
		t.Error("Failed to unlink file, error : ", code)
	}

	dirEntries, code = pfs.OpenDir("", nil)
	if code != fuse.OK {
		t.Error("Could not open directory, error : ", code)
	}
	if len(dirEntries) != 0 {
		t.Error("Incorrect number of files in directory : ", len(dirEntries))
	}
}

func TestFuseLink(t *testing.T) {
	setuptesting(t)
	defer removeTestDir("pfiTestPfsDir")

	_, err := commands.CreatCommand(path.Join(os.TempDir(), "pfiTestPfsDir"), "helloworld.txt", os.FileMode(0777))
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}

	pfs := &ParanoidFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
	}

	_, code := pfs.Create("file1", 0, uint32(os.FileMode(0777)), nil)
	if code != fuse.OK {
		t.Error("Create did not return OK. Actual : ", code)
	}

	code = pfs.Link("file1", "file2", nil)
	if code != fuse.OK {
		t.Error("Link did not return OK. Actual : ", code)
	}

	file, code := pfs.Open("file2", uint32(os.O_RDWR), nil)
	if code != fuse.OK {
		t.Error("Did not get OK opening file. Actual :", code)
	}

	_, code = file.Write([]byte("testhello"), 0)
	if code != fuse.OK {
		t.Error("Write did not return OK. Actual :", code)
	}

	buf := make([]byte, 9)
	_, code = file.Read(buf, 0)
	if code != fuse.OK {
		t.Error("Read did not return OK. Actual :", code)
	}

	if string(buf) != "testhello" {
		t.Error("Read did not return correct result. Actual :", string(buf))
	}
}

func TestFuseUtimes(t *testing.T) {
	setuptesting(t)
	defer removeTestDir("pfiTestPfsDir")

	pfs := &ParanoidFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
	}

	file, code := pfs.Create("helloworld.txt", 0, uint32(os.FileMode(0777)), nil)
	if code != fuse.OK {
		t.Error("Failed to create file, error : ", code)
	}

	atime := time.Unix(100, 101*1000)
	mtime := time.Unix(500, 530*1000)
	roundFactor := time.Duration(1 * time.Second)
	code = file.Utimens(&atime, &mtime)
	if code != fuse.OK {
		t.Error("Failed to utimens file, error : ", code)
	}

	attr, code := pfs.GetAttr("helloworld.txt", nil)
	if code != fuse.OK {
		t.Error("Failed to stat file, error : ", code)
	}
	if attr.ModTime().Round(roundFactor) != mtime.Round(roundFactor) {
		t.Error("Incorrect mtime received : ", attr.ModTime().Round(roundFactor))
	}
	if attr.AccessTime().Round(roundFactor) != atime.Round(roundFactor) {
		t.Error("Incorrect atime recieved : ", attr.AccessTime().Round(roundFactor))
	}
}
