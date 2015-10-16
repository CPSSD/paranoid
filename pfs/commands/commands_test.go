package commands

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
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

func doWriteCommand(t *testing.T, file, data string, offset, length int) {
	cmd := exec.Command("go", "run", "../main.go", "write", path.Join(os.TempDir(), "paranoidTest"), file)
	if offset != -1 {
		if length != -1 {
			cmd = exec.Command("go", "run", "../main.go", "write", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset), strconv.Itoa(length))
		} else {
			cmd = exec.Command("go", "run", "../main.go", "write", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset))
		}
	}
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

func doReadCommand(t *testing.T, file string, offset, length int) string {
	cmd := exec.Command("go", "run", "../main.go", "read", path.Join(os.TempDir(), "paranoidTest"), file)
	if offset != -1 {
		if length != -1 {
			cmd = exec.Command("go", "run", "../main.go", "read", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset), strconv.Itoa(length))
		} else {
			cmd = exec.Command("go", "run", "../main.go", "read", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset))
		}
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		t.Error("Error running read command: ", err)
	}
	return out.String()
}

func doReadDirCommand(t *testing.T) []string {
	cmd := exec.Command("go", "run", "../main.go", "readdir", path.Join(os.TempDir(), "paranoidTest"))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		t.Error("Error running readdir command: ", err)
	}
	anwser := strings.Split(out.String(), "\n")
	return anwser[0 : len(anwser)-1]
}

func TestSimpleCommandUsage(t *testing.T) {
	createTestDir(t)
	defer removeTestDir()
	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt"}
	CreatCommand(args)
	doWriteCommand(t, "test.txt", "BLAH #1", -1, -1)
	returnData := strings.TrimRight(doReadCommand(t, "test.txt", -1, -1), "\x00")
	if returnData != "BLAH #1" {
		t.Error("Output does not match ", returnData)
	}
}

func TestComplexCommandUsage(t *testing.T) {
	createTestDir(t)
	defer removeTestDir()
	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt"}
	CreatCommand(args)
	doWriteCommand(t, "test.txt", "START", -1, -1)
	returnData := strings.TrimRight(doReadCommand(t, "test.txt", 2, 2), "\x00")
	if returnData != "AR" {
		t.Error("Output from partial read does not match")
	}
	doWriteCommand(t, "test.txt", "END", 5, -1)
	returnData = strings.TrimRight(doReadCommand(t, "test.txt", -1, -1), "\x00")
	if returnData != "STARTEND" {
		t.Error("Output from full read does not match")
	}
	files := doReadDirCommand(t)
	if files[0] != "test.txt" || len(files) > 1 {
		t.Error("Readdir got incorrect result")
	}
}
