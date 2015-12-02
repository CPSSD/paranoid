// +build integration

package main

import (
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfi/filesystem"
	"github.com/cpssd/paranoid/pfi/pfsminterface"
	"github.com/cpssd/paranoid/pfi/util"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"os"
	"os/exec"
	"path"
	"strings"
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
	util.Log = logger.New("testPackage", "testComponent", "/dev/null")
	os.Exit(m.Run())
}

func TestFuseExternalUsage(t *testing.T) {
	createTestDir(t, "pfiTestPfsDir")
	defer removeTestDir("pfiTestPfsDir")
	createTestDir(t, "pfiTestMountPoint")
	defer removeTestDir("pfiTestMountPoint")
	cmd := exec.Command("pfsm", "init", path.Join(os.TempDir(), "pfiTestPfsDir"))
	err := cmd.Run()
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}
	cmd = exec.Command("go", "run", "main.go", "-n", path.Join(os.TempDir(), "pfiTestPfsDir"), path.Join(os.TempDir(), "pfiTestMountPoint"))
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		t.Error("Unable to start pfi :", err)
	}
	time.Sleep(time.Second * 2) //Wait 2 seconds so that file system is mounted before we start testing it.

	cmd = exec.Command("touch", path.Join(os.TempDir(), "pfiTestMountPoint", "helloworld.txt"))
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Error("Error running touch :", err)
	}

	cmd = exec.Command("ls", path.Join(os.TempDir(), "pfiTestMountPoint"))
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		t.Error("Error running ls : ", err)
	}
	fileName := strings.TrimSpace(string(output))
	if fileName != "helloworld.txt" {
		t.Error("ls returned incorrect result")
	}

	cmd = exec.Command("pfsm", "-n", "write", path.Join(os.TempDir(), "pfiTestPfsDir"), "helloworld.txt")
	cmd.Stderr = os.Stderr
	pipe, err := cmd.StdinPipe()
	if err != nil {
		t.Error("Error creating pipe :", err)
	}
	_, err = pipe.Write([]byte("CPSSD 4 Lyfe"))
	if err != nil {
		t.Error("Error writing to pipe :", err)
	}
	err = pipe.Close()
	if err != nil {
		t.Error("Error closing pipe :", err)
	}
	err = cmd.Run()
	if err != nil {
		t.Error("Error writing to pfsm :", err)
	}

	time.Sleep(time.Second * 1) //Wait before cating or old data may be recieved
	cmd = exec.Command("cat", path.Join(os.TempDir(), "pfiTestMountPoint", "helloworld.txt"))
	cmd.Stderr = os.Stderr
	output, err = cmd.Output()
	if err != nil {
		t.Error("Error running cat command :", err)
	}
	if string(output) != "CPSSD 4 Lyfe" {
		t.Error("Unexpected output from cat command")
	}

	cmd = exec.Command("stat", path.Join(os.TempDir(), "pfiTestMountPoint", "helloworld.txt"))
	cmd.Stderr = os.Stderr
	output, err = cmd.Output()
	if err != nil || len(output) == 0 {
		t.Error("Error running stat command :", err)
	}

	cmd = exec.Command("fusermount", "-u", "-z", path.Join(os.TempDir(), "pfiTestMountPoint"))
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Error("could not dismount filesystem", err)
	}
}

func TestFuseFilePerms(t *testing.T) {
	createTestDir(t, "pfiTestPfsDir")
	defer removeTestDir("pfiTestPfsDir")
	util.PfsDirectory = path.Join(os.TempDir(), "pfiTestPfsDir")
	pfsminterface.OriginFlag = "-n"

	cmd := exec.Command("pfsm", "init", path.Join(os.TempDir(), "pfiTestPfsDir"))
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}

	cmd = exec.Command("pfsm", "-n", "creat", path.Join(os.TempDir(), "pfiTestPfsDir"), "helloworld.txt", "0777")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}

	pfs := &filesystem.ParanoidFileSystem{
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
	createTestDir(t, "pfiTestPfsDir")
	defer removeTestDir("pfiTestPfsDir")
	util.PfsDirectory = path.Join(os.TempDir(), "pfiTestPfsDir")
	pfsminterface.OriginFlag = "-n"

	cmd := exec.Command("pfsm", "init", path.Join(os.TempDir(), "pfiTestPfsDir"))
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}

	pfs := &filesystem.ParanoidFileSystem{
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
	createTestDir(t, "pfiTestPfsDir")
	defer removeTestDir("pfiTestPfsDir")
	util.PfsDirectory = path.Join(os.TempDir(), "pfiTestPfsDir")
	pfsminterface.OriginFlag = "-n"

	cmd := exec.Command("pfsm", "init", path.Join(os.TempDir(), "pfiTestPfsDir"))
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}

	pfs := &filesystem.ParanoidFileSystem{
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
	createTestDir(t, "pfiTestPfsDir")
	defer removeTestDir("pfiTestPfsDir")
	util.PfsDirectory = path.Join(os.TempDir(), "pfiTestPfsDir")
	pfsminterface.OriginFlag = "-n"

	cmd := exec.Command("pfsm", "init", path.Join(os.TempDir(), "pfiTestPfsDir"))
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}

	pfs := &filesystem.ParanoidFileSystem{
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
	createTestDir(t, "pfiTestPfsDir")
	defer removeTestDir("pfiTestPfsDir")
	util.PfsDirectory = path.Join(os.TempDir(), "pfiTestPfsDir")
	pfsminterface.OriginFlag = "-n"

	cmd := exec.Command("pfsm", "init", path.Join(os.TempDir(), "pfiTestPfsDir"))
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Error("pfsm setup failed :", err)
	}

	pfs := &filesystem.ParanoidFileSystem{
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
