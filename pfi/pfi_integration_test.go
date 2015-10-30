// +build integration

package main

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

func createTestDir(t *testing.T, name string) {
	err := os.Mkdir(path.Join(os.TempDir(), name), 0777)
	if err != nil {
		t.Error(err)
	}
}

func removeTestDir(name string) {
	os.RemoveAll(path.Join(os.TempDir(), name))
}

func TestFuseUsage(t *testing.T) {
	createTestDir(t, "pfiTestPfsDir")
	defer removeTestDir("pfiTestPfsDir")
	createTestDir(t, "pfiTestMountPoint")
	defer removeTestDir("pfiTestMountPoint")
	cmd := exec.Command("pfs", "init", path.Join(os.TempDir(), "pfiTestPfsDir"))
	err := cmd.Run()
	if err != nil {
		t.Error("Pfs setup failed :", err)
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

	cmd = exec.Command("pfs", "-n", "write", path.Join(os.TempDir(), "pfiTestPfsDir"), "helloworld.txt")
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
		t.Error("Error writing to pfs :", err)
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

	cmd = exec.Command("fusermount", "-u", "-z", path.Join(os.TempDir(), "pfiTestMountPoint"))
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Error("could not dismount filesystem", err)
	}
}
