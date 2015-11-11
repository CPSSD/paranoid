// +build integration

package main

import (
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
	createTestDir(t, "pfiTestMountPoint")
	defer removeTestDir("pfiTestMountPoint")
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
