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
	"time"
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

func runCommand(t *testing.T, stdinData []byte, cmdArgs ...string) (int, string) {
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
			code, _ := runCommand(t, []byte(data), "write", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset), strconv.Itoa(length))
			return code
		} else {
			code, _ := runCommand(t, []byte(data), "write", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset))
			return code
		}
	}
	code, _ := runCommand(t, []byte(data), "write", path.Join(os.TempDir(), "paranoidTest"), file)
	return code
}

func doReadCommand(t *testing.T, file string, offset, length int) (int, string) {
	if offset != -1 {
		if length != -1 {
			return runCommand(t, nil, "read", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset), strconv.Itoa(length))
		} else {
			return runCommand(t, nil, "read", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(offset))
		}
	}
	return runCommand(t, nil, "read", path.Join(os.TempDir(), "paranoidTest"), file)
}

func doReadDirCommand(t *testing.T) (int, []string) {
	code, data := runCommand(t, nil, "readdir", path.Join(os.TempDir(), "paranoidTest"))
	anwser := strings.Split(data, "\n")
	return code, anwser[0 : len(anwser)-1]
}

func doStatCommand(t *testing.T, file string) (int, statInfo) {
	code, data := runCommand(t, nil, "stat", path.Join(os.TempDir(), "paranoidTest"), file)
	stats := statInfo{}
	if code != returncodes.OK {
		return code, stats
	}
	err := json.Unmarshal([]byte(data), &stats)
	if err != nil {
		t.Error("Error parsing stat data = ", data)
	}
	return code, stats
}

func doUtimesCommand(t *testing.T, file string, atime, mtime *time.Time) int {
	timeStuct := &timeInfo{
		Atime: atime,
		Mtime: mtime,
	}
	data, err := json.Marshal(timeStuct)
	if err != nil {
		t.Error("Error marshalling utimes data", err)
	}
	code, _ := runCommand(t, data, "utimes", path.Join(os.TempDir(), "paranoidTest"), file)
	return code
}

func doAccessCommand(t *testing.T, file string, mode int) int {
	code, _ := runCommand(t, nil, "access", path.Join(os.TempDir(), "paranoidTest"), file, strconv.Itoa(mode))
	return code
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
		t.Error("Write did not return OK. Actual:", code)
	}
	code, returnData := doReadCommand(t, "test.txt", -1, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK. Actual:", code)
	}
	if returnData != "BLAH #1" {
		t.Error("Output does not match BLAH #1. Actual:", returnData)
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
		t.Error("Read did not return OK. Actual:", code)
	}
	code, returnData := doReadCommand(t, "test.txt", 2, 2)
	if code != returncodes.OK {
		t.Error("Read did not return OK. Actual:", code)
	}
	if returnData != "AR" {
		t.Error("Output from partial read does not match ", returnData)
	}
	code = doWriteCommand(t, "test.txt", "END", 5, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK ", code)
	}
	code, returnData = doReadCommand(t, "test.txt", -1, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK. Actual:", code)
	}
	if returnData != "STARTEND" {
		t.Error("Output from full read does not match STARTEND. Actual:", returnData)
	}
	code, files := doReadDirCommand(t)
	if code != returncodes.OK {
		t.Error("Read did not return OK. Actual:", code)
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
		t.Error("Stat did not return OK. Actual:", code)
	}
	if statIn.Perms != 0777 {
		t.Error("Incorrect file permissions = ", statIn.Perms)
	}

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "377"}
	ChmodCommand(args)

	code, statIn = doStatCommand(t, "test.txt")
	if code != returncodes.OK {
		t.Error("Stat did not return OK. Actual:", code)
	}
	if statIn.Perms != 0377 {
		t.Error("Incorrect file permissions = ", statIn.Perms)
	}

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "500"}
	ChmodCommand(args)

	code, statIn = doStatCommand(t, "test.txt")
	if code != returncodes.OK {
		t.Error("Stat did not return OK. Actual:", code)
	}
	if statIn.Perms != 0500 {
		t.Error("Incorrect file permissions = ", statIn.Perms)
	}

	code = doWriteCommand(t, "test.txt", "helloWorld", -1, -1)
	if code != returncodes.EACCES {
		t.Error("Should not be able to write file ", statIn.Perms)
	}

	code = doAccessCommand(t, "test.txt", 4)
	if code != returncodes.OK {
		t.Error("Access command did not return OK. Actual:", code)
	}
	code = doAccessCommand(t, "test.txt", 2)
	if code != returncodes.EACCES {
		t.Error("Access command did not return EACCES. Actual:", code)
	}
}

func TestENOENT(t *testing.T) {
	Flags.Network = true
	createTestDir(t)
	defer removeTestDir()

	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)

	code, _ := doReadCommand(t, "test.txt", -1, -1)
	if code != returncodes.ENOENT {
		t.Error("Read did not return ENOENT. Actual:", code)
	}

	code, _ = doStatCommand(t, "test.txt")
	if code != returncodes.ENOENT {
		t.Error("Stat did not return ENOENT. Actual:", code)
	}

	code = doWriteCommand(t, "test.txt", "data", -1, -1)
	if code != returncodes.ENOENT {
		t.Error("Write did not return ENOENT. Actual:", code)
	}
}

func TestFilesystemCommands(t *testing.T) {
	Flags.Network = true
	createTestDir(t)
	defer removeTestDir()

	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "777"}
	CreatCommand(args)

	code, files := doReadDirCommand(t)
	if code != returncodes.OK {
		t.Error("Readdir did not return OK. Actual:", code)
	}
	if files[0] != "test.txt" || len(files) > 1 {
		t.Error("Readdir got incorrect result")
	}

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "test2.txt"}
	RenameCommand(args)

	code, files = doReadDirCommand(t)
	if code != returncodes.OK {
		t.Error("Readdir did not return OK. Actual:", code)
	}
	if files[0] != "test2.txt" || len(files) > 1 {
		t.Error("Readdir got incorrect result")
	}

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test2.txt"}
	UnlinkCommand(args)

	code, files = doReadDirCommand(t)
	if code != returncodes.OK {
		t.Error("Readdir did not return OK. Actual:", code)
	}
	if len(files) > 0 {
		t.Error("Readdir got incorrect result")
	}
}

func TestLinkCommand(t *testing.T) {
	Flags.Network = true
	createTestDir(t)
	defer removeTestDir()

	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "777"}
	CreatCommand(args)

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "test2.txt"}
	LinkCommand(args)

	code, files := doReadDirCommand(t)
	if code != returncodes.OK {
		t.Error("Readdir did not return OK. Actual:", code)
	}
	if files[0] != "test.txt" && files[1] != "test.txt" {
		t.Error("Readdir got incorrect result")
	}
	if files[0] != "test2.txt" && files[1] != "test2.txt" {
		t.Error("Readdir got incorrect result")
	}
	if len(files) != 2 {
		t.Error("Readdir got incorrect results")
	}

	code = doWriteCommand(t, "test2.txt", "hellotest", -1, -1)
	if code != returncodes.OK {
		t.Error("Write did not return OK. Actual:", code)
	}

	code, data := doReadCommand(t, "test.txt", -1, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK. Actual:", code)
	}
	if data != "hellotest" {
		t.Error("Read did not return correct result")
	}

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt"}
	UnlinkCommand(args)

	code, files = doReadDirCommand(t)
	if code != returncodes.OK {
		t.Error("Readdir did not return OK. Actual:", code)
	}
	if files[0] != "test2.txt" {
		t.Error("Readdir got incorrect result")
	}
	if len(files) != 1 {
		t.Error("Readdir got incorrect result")
	}

	code, data = doReadCommand(t, "test2.txt", -1, -1)
	if code != returncodes.OK {
		t.Error("Read did not return OK. Actual:", code)
	}
	if data != "hellotest" {
		t.Error("Read did not return correct result")
	}
}

func TestUtimes(t *testing.T) {
	Flags.Network = true
	createTestDir(t)
	defer removeTestDir()

	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "777"}
	CreatCommand(args)

	atime := time.Unix(100, 100)
	mtime := time.Unix(500, 250)
	code := doUtimesCommand(t, "test.txt", &atime, &mtime)

	code, statIn := doStatCommand(t, "test.txt")
	if code != returncodes.OK {
		t.Error("Stat did not return OK. Actual:", code)
	}
	if statIn.Atime != time.Unix(100, 100) {
		t.Error("Incorrect stat time. Acutal: ", statIn.Atime)
	}
	if statIn.Mtime != time.Unix(500, 250) {
		t.Error("Incorrect stat time. Acutal: ", statIn.Atime)
	}
}

func TestTruncate(t *testing.T) {
	Flags.Network = true
	createTestDir(t)
	defer removeTestDir()

	args := []string{path.Join(os.TempDir(), "paranoidTest")}
	InitCommand(args)
	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "777"}
	CreatCommand(args)

	code := doWriteCommand(t, "test.txt", "HI!!!!!", -1, -1)
	if code != returncodes.OK {
		t.Error("Write command failed!")
	}

	args = []string{path.Join(os.TempDir(), "paranoidTest"), "test.txt", "3"}
	TruncateCommand(args)

	code, data := doReadCommand(t, "test.txt", -1, -1)
	if code != returncodes.OK {
		t.Error("Read command did not return OK. Actual:", code)
	}
	if data != "HI!" {
		t.Error("Read command returned incorrect output ", data)
	}
}
