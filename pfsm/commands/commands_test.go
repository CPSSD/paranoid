// +build !integration

package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"
)

func createTestDir(t *testing.T) {
	err := os.RemoveAll(path.Join(os.TempDir(), "paranoidTest"))
	if err != nil {
		t.Error(err)
	}
	err = os.Mkdir(path.Join(os.TempDir(), "paranoidTest"), 0777)
	if err != nil {
		t.Error(err)
	}
}

func removeTestDir() {
	os.RemoveAll(path.Join(os.TempDir(), "paranoidTest"))
}

func RunCommand(t *testing.T, stdinData []byte, cmdArgs ...string) (int, string) {
	cmdArgs = append(cmdArgs, "-n")
	command := exec.Command("pfsm", cmdArgs...)
	command.Stderr = os.Stderr

	if stdinData != nil {
		stdinPipe, err := command.StdinPipe()
		if err != nil {
			t.Error("Error running pfsm "+cmdArgs[0]+" :", err)
		}
		_, err = stdinPipe.Write(stdinData)
		if err != nil {
			t.Error("Error running pfsm "+cmdArgs[0]+" :", err)
		}
		err = stdinPipe.Close()
		if err != nil {
			t.Error("Error running pfsm "+cmdArgs[0]+" :", err)
		}
	}

	output, err := command.Output()
	if err != nil {
		t.Error("Error running pfsm "+cmdArgs[0]+" :", err)
	}
	code, err := strconv.Atoi(string(output[0:2]))
	if err != nil {
		t.Error("Error running pfsm "+cmdArgs[0]+" (invalid return code) :", err)
	}
	return code, string(output[2:])
}

func doWriteCommand(t *testing.T, file, data string, offset, length int) int {
	if offset != -1 {
		if length != -1 {
			code, _ := RunCommand(t, []byte(data), "write", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset), strconv.Itoa(length))
			return code
		} else {
			code, _ := RunCommand(t, []byte(data), "write", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset))
			return code
		}
	}
	code, _ := RunCommand(t, []byte(data), "write", path.Join(os.TempDir(), "paranoidTest"), file)
	return code
}

func doReadCommand(t *testing.T, file string, offset, length int) (int, string) {
	if offset != -1 {
		if length != -1 {
			return RunCommand(t, nil, "read", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset), strconv.Itoa(length))
		} else {
			return RunCommand(t, nil, "read", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset))
		}
	}
	return RunCommand(t, nil, "read", path.Join(os.TempDir(), "paranoidTest"), file)
}

func doReadDirCommand(t *testing.T) (int, []string) {
	code, data := RunCommand(t, nil, "readdir", path.Join(os.TempDir(), "paranoidTest"))
	anwser := strings.Split(data, "\n")
	return code, anwser[0 : len(anwser)-1]
}

func doStatCommand(t *testing.T, file string) (int, statInfo) {
	code, data := RunCommand(t, nil, "stat", path.Join(os.TempDir(), "paranoidTest"), file)
	stats := statInfo{}
	err := json.Unmarshal([]byte(data), &stats)
	if err != nil {
		t.Error("Error parsing stat data = ", data)
	}
	return code, stats
}

func TestSimpleCommandUsage(t *testing.T) {
	Flags.Network = true // so that the tests don't try to make a network connection
	createTestDir(t)
	defer removeTestDir()
	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "777"}
	CreatCommand(args)
	code := doWriteCommand(t, "test.txt", "BLAH #1", -1, -1)
	if code != returncodes.OK {
		t.Error("Write did not return OK")
	}
	code, returnData := doReadCommand(t, "test.txt", -1, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK")
	}
	if returnData != "BLAH #1" {
		t.Error("Output does not match ", returnData)
	}
}

func TestComplexCommandUsage(t *testing.T) {
	Flags.Network = true
	createTestDir(t)
	defer removeTestDir()
	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "777"}
	CreatCommand(args)
	code := doWriteCommand(t, "test.txt", "START", -1, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK")
	}
	code, returnData := doReadCommand(t, "test.txt", 2, 2)
	if code != returncodes.OK {
		t.Error("Read did not return OK")
	}
	if returnData != "AR" {
		t.Error("Output from partial read does not match")
	}
	code = doWriteCommand(t, "test.txt", "END", 5, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK")
	}
	code, returnData = doReadCommand(t, "test.txt", -1, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK")
	}
	if returnData != "STARTEND" {
		t.Error("Output from full read does not match")
	}
	code, files := doReadDirCommand(t)
	if code != returncodes.OK {
		t.Error("Read did not return OK")
	}
	if files[0] != "test.txt" || len(files) > 1 {
		t.Error("Readdir got incorrect result")
	}
}

func TestFilePermissionsCommands(t *testing.T) {
	Flags.Network = true
	createTestDir(t)
	defer removeTestDir()
	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "777"}
	CreatCommand(args)

	code, statIn := doStatCommand(t, "test.txt")
	if code != returncodes.OK {
		t.Error("Stat did not return OK")
	}
	if statIn.Perms != 0777 {
		t.Error("Incorrect file permissions = ", statIn.Perms)
	}

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "377"}
	ChmodCommand(args)

	code, statIn = doStatCommand(t, "test.txt")
	if code != returncodes.OK {
		t.Error("Stat did not return OK")
	}
	if statIn.Perms != 0377 {
		t.Error("Incorrect file permissions = ", statIn.Perms)
	}

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "500"}
	ChmodCommand(args)

	code, statIn = doStatCommand(t, "test.txt")
	if code != returncodes.OK {
		t.Error("Stat did not return OK")
	}
	if statIn.Perms != 0500 {
		t.Error("Incorrect file permissions = ", statIn.Perms)
	}

	code = doWriteCommand(t, "test.txt", "helloWorld", -1, -1)
	if code != returncodes.EACCES {
		t.Error("Should not be able to write file ", statIn.Perms)
	}
}
