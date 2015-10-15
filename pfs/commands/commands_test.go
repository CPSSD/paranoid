package commands

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func createTestDir(t *testing.T) {
	err := os.Mkdir(path.Join(os.TempDir(), "paranoidTest"), 0777)
	if err != nil {
		t.Error(err)
	}
}

func removeTestDir() {
	os.RemoveAll(path.Join(os.TempDir(), "paranoidTest"))
}

func doWriteCommand(t *testing.T, file, data string) {
	cmd := exec.Command("go", "run", "../main.go", "write", path.Join(os.TempDir(), "paranoidTest"), file)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Error("Error setting up write command")
	}
	io.WriteString(stdin, data)
	stdin.Close()
	err = cmd.Run()
	if err != nil {
		t.Error("write command could not start")
	}
}

func doReadCommand(t *testing.T, file string) string {
	cmd := exec.Command("go", "run", "../main.go", "read", path.Join(os.TempDir(), "paranoidTest"), file)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		t.Error("Error running read command: ", err)
	}
	return out.String()
}

func TestSimpleCommandUsage(t *testing.T) {
	createTestDir(t)
	defer removeTestDir()
	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt"}
	CreatCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt"}
	doWriteCommand(t, "test.txt", "BLAH #1")
	returnData := strings.TrimRight(doReadCommand(t, "test.txt"), "\x00")
	if returnData != "BLAH #1" {
		t.Error("Output does not match ", returnData)
	}
}
