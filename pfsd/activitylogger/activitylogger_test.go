// +build !integration

package activitylogger

import (
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/activitylogger"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var (
	Log     *logger.ParanoidLogger
	testDir string
)

func createTestDir() {
	err := os.RemoveAll(testDir)
	if err != nil {
		Log.Fatal("error creating test directory:", err)
	}

	err = os.Mkdir(testDir, 0777)
	if err != nil {
		Log.Fatal("error creating test directory:", err)
	}

	err = os.Mkdir(path.Join(testDir, "meta"), 0777)
	if err != nil {
		Log.Fatal("error creating test directory:", err)
	}
}

func removeTestDir() {
	os.RemoveAll(testDir)
}

func TestWriteReadDelete(t *testing.T) {
	Log = logger.New("activitylogger_test", "pfsdTest", os.DevNull)
	testDir = path.Join(os.TempDir(), "paranoidTest")
	removeTestDir()
	createTestDir()

	al := New(testDir)
	logDir := path.Join(testDir, "meta", "activity_logs")

	// Testing Write functionality
	i, err := al.WriteEntry(&pb.Entry{
		Type: 0,
		Path: "ThisIsAPath",
	})
	if err != nil || i != 1 {
		t.Error("received error writing log, err:", err)
	}

	i, err = al.WriteEntry(&pb.Entry{
		Type: 0,
		Path: "ThisIsAPath2",
	})
	if err != nil || i != 2 {
		t.Error("received error writing log, err:", err)
	}

	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		t.Error("Error reading log directory")
	}

	if len(files) != 2 {
		t.Error("number of files in directory is not what it shoudl be, writing error.")
	}
	if files[0].Name() != "1000001" || files[1].Name() != "1000002" {
		t.Error("Files not named what they should be, file1: ", files[0].Name(), "file2: ", files[1].Name())
	}

	// Testing Read functionality
	li := al.LastEntryIndex()
	if li != 2 {
		t.Error("LastEntryIndex not what it should be: ", li)
	}

	e, err := al.GetEntry(al.LastEntryIndex())
	if err != nil {
		t.Error("Error received when reading log: ", err)
	}

	if e.Type != 0 || e.Path != "ThisIsAPath2" {
		t.Error("Bad protobuf received from read: ", e)
	}

	// Testing Delete functionality
	err = al.DeleteEntry(1)
	if err != nil {
		t.Error("Error received when deleting log: ", err)
	}

	files, err = ioutil.ReadDir(logDir)
	if err != nil {
		t.Error("Error reading log directory")
	}

	if len(files) != 0 {
		t.Error("number of files in directory is not what it shoudl be, delete error.")
	}
}
