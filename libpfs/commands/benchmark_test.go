package commands

import (
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

func BenchmarkCreat(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
	}
}

func BenchmarkWrite(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		WriteCommand(testDirectory, "test.txt"+str, 0, 0, []byte("Hello World"), false)
	}
}

func BenchmarkRename(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		RenameCommand(testDirectory, "test.txt"+str, "test2.txt"+str, false)
	}
}

func BenchmarkRead(b *testing.B) {
	setupTestDirectory()
	CreatCommand(testDirectory, "test.txt", os.FileMode(0777), false)
	for n := 0; n < b.N; n++ {
		ReadCommand(testDirectory, "test.txt", 0, 0)
	}
}

func BenchmarkStat(b *testing.B) {
	setupTestDirectory()
	CreatCommand(testDirectory, "test.txt", os.FileMode(0777), false)
	for n := 0; n < b.N; n++ {
		StatCommand(testDirectory, "test.txt")
	}
}

func BenchmarkTruncate(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		TruncateCommand(testDirectory, "test.txt"+str, 3, false)
	}
}

func BenchmarkUtimes(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		atime := time.Unix(100, 100)
		mtime := time.Unix(500, 250)
		CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		UtimesCommand(testDirectory, "test.txt"+str, &atime, &mtime, false)
	}
}

func BenchmarkMkDir(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		MkdirCommand(testDirectory, "testDir"+str, os.FileMode(0777), false)
	}
}

func BenchmarkReadDir(b *testing.B) {
	setupTestDirectory()
	MkdirCommand(testDirectory, "testDir", os.FileMode(0777), false)
	CreatCommand(path.Join(testDirectory, "testDir"), "test.txt", os.FileMode(0777), false)
	for n := 0; n < b.N; n++ {
		ReadDirCommand(testDirectory, "testDir")
	}
}

func BenchmarkLink(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		LinkCommand(testDirectory, "test.txt"+str, "test2.txt"+str, false)
	}
}

func BenchmarkSymLink(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		SymlinkCommand(testDirectory, "testfolder/testlink", "testsymlink"+str, false)
	}
}

func BenchmarkReadLink(b *testing.B) {
	setupTestDirectory()
	SymlinkCommand(testDirectory, "testfolder/testlink", "testsymlink", false)
	for n := 0; n < b.N; n++ {
		ReadlinkCommand(testDirectory, "testsymlink")
	}
}
