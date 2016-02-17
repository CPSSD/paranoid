// +build !integration !benchmark

package logger

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestOutput(t *testing.T) {
	log := New("testPackage", "testComponent", os.DevNull)
	err := log.SetOutput(STDERR)
	if err != nil {
		t.Fatal("Failed to set logger output:", err)
	}

	const testString = "test"

	var b bytes.Buffer
	log.AddAdditionalWriter(&b)
	log.Info(testString)

	expected := "[INFO]  testPackage: " + testString + "\n"

	result := (b.String())[20:] // Ignore the time/date

	if result != expected {
		t.Errorf("Expected %s, got %s\n", expected, result)
	}
}

func TestOutputf(t *testing.T) {
	log := New("testPackage", "testComponent", "/dev/null")
	err := log.SetOutput(STDERR)
	if err != nil {
		t.Fatal("Failed to set logger output:", err)
	}

	testArgs := []string{"testy %s", "test"}

	var b bytes.Buffer
	log.AddAdditionalWriter(&b)
	log.Infof(testArgs[0], testArgs[1])

	expected := "[INFO]  testPackage: testy test\n"

	result := (b.String())[20:] // Ignore the time/date

	if result != expected {
		t.Errorf("Expected \"%s\", got \"%s\"\n", expected, result)
	}
}

func TestLogLevel(t *testing.T) {
	log := New("testPackage", "testComponent", "/dev/null")
	err := log.SetOutput(STDERR)
	if err != nil {
		t.Fatal("Failed to set logger output:", err)
	}
	log.SetLogLevel(INFO)

	var b bytes.Buffer
	log.AddAdditionalWriter(&b)
	log.Debug("test")

	if len(b.String()) != 0 {
		t.Errorf("%s returned. Expected nothing", b.String())
	}
}

func TestLogFile(t *testing.T) {
	os.Mkdir("/tmp/pfsLogTest", 0777)
	// Remove the file that the logger is saving to after testing
	defer os.RemoveAll("/tmp/pfsLogTest")
	log := New("testPackage", "testComponent", "/tmp/pfsLogTest")
	err := log.SetOutput(STDERR | LOGFILE)
	if err != nil {
		t.Fatal("Failed to set logger output:", err)
	}

	const testString = "test"
	expected := "[INFO]  testPackage: " + testString + "\n"

	log.Info(testString)
	data, _ := ioutil.ReadFile("/tmp/pfsLogTest/testComponent.log")
	if len(data) == 0 && string(data[20:]) != expected {
		t.Errorf("Expected %s, got %s\n", expected, string(data[20:]))
	}
}
