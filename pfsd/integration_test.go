// +build integration

package main

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func createTestDir(t *testing.T, name string) {
	os.RemoveAll(path.Join(os.TempDir(), name))
	err := os.Mkdir(path.Join(os.TempDir(), name), 0777)
	if err != nil {
		t.Fatal("Error creating directory", err)
	}
}

func removeTestDir(name string) {
	time.Sleep(1 * time.Second)
	os.RemoveAll(path.Join(os.TempDir(), name))
}

func TestKillSignal(t *testing.T) {
	createTestDir(t, "testksMountpoint")
	defer removeTestDir("testksMountpoint")
	createTestDir(t, "testksDirectory")
	defer removeTestDir("testksDirectory")

	commands.Log = logger.New("pfsdintegration", "pfsdintegration", os.DevNull)

	discovery := exec.Command("discovery-server", "--port=10102")
	err := discovery.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := discovery.Process.Kill()
		if err != nil {
			t.Error("Failed to kill discovery server,", err)
		}
	}()

	_, err = commands.InitCommand(path.Join(os.TempDir(), "testksDirectory"))
	if err != nil {
		t.Fatal(err)
	}

	pfsd := exec.Command("pfsd", path.Join(os.TempDir(), "testksDirectory"), path.Join(os.TempDir(), "testksMountpoint"), "localhost", "10102", "testPool")
	err = pfsd.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer pfsd.Process.Kill()
	defer func() {
		cmd := exec.Command("fuserunmount", "-z", "-u", path.Join(os.TempDir(), "testksMountPoint"))
		cmd.Run()
	}()

	time.Sleep(5 * time.Second)

	pidPath := path.Join(os.TempDir(), "testksDirectory", "meta", "pfsd.pid")
	if _, err := os.Stat(pidPath); err == nil {
		pidByte, err := ioutil.ReadFile(pidPath)
		if err != nil {
			t.Fatal("Can't read pid file", err)
		}
		pid, err := strconv.Atoi(string(pidByte))
		if err != nil {
			t.Fatal("Incorrect pid information", err)
		}
		err = syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			t.Fatal("Can not kill PFSD,", err)
		}

		done := make(chan bool, 1)
		go func() {
			pfsd.Wait()
			done <- true
		}()

		select {
		case <-time.After(10 * time.Second):
			t.Fatal("pfsd did not finish within 10 seconds")
		case <-done:
			break
		}
	} else {
		t.Fatal("Could not read pid file:", err)
	}
}
