package benchmark

import (
	. "github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/logger"
	"log"
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

var testDirectory string

func TestMain(m *testing.M) {
	Log = logger.New("commandsTest", "pfsmTest", os.DevNull)
	Log.SetLogLevel(logger.ERROR)
	testDirectory = path.Join(os.TempDir(), "paranoidTest")
	defer removeTestDir()
	os.Exit(m.Run())
}

func removeTestDir() {
	os.RemoveAll(testDirectory)
}

func createTestDir() {
	err := os.RemoveAll(testDirectory)
	if err != nil {
		Log.Fatal("error creating test directory:", err)
	}

	err = os.Mkdir(testDirectory, 0777)
	if err != nil {
		Log.Fatal("error creating test directory:", err)
	}
}

func setupTestDirectory() {
	createTestDir()

	code, err := InitCommand(testDirectory)
	if code != returncodes.OK {
		Log.Fatal("error initing directory for testing:", err)
	}
}

func BenchmarkCreat(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		code, err := CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		if code != returncodes.OK {
			log.Fatal("error creating test file:", err)
		}
	}
}

func BenchmarkWrite(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		code, err := CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		if code != returncodes.OK {
			log.Fatal("error creating test file:", err)
		}
		code, err, _ = WriteCommand(testDirectory, "test.txt"+str, 0, 0, []byte("Hello World"), false)
		if code != returncodes.OK {
			log.Fatal("error writing to test file:", err)
		}
	}
}

func BenchmarkRename(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		code, err := CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		if code != returncodes.OK {
			log.Fatal("error creating test file:", err)
		}
		code, err = RenameCommand(testDirectory, "test.txt"+str, "test2.txt"+str, false)
		if code != returncodes.OK {
			log.Fatal("error renaming test file:", err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	setupTestDirectory()
	code, err := CreatCommand(testDirectory, "test.txt", os.FileMode(0777), false)
	if code != returncodes.OK {
		log.Fatal("error creating test file:", err)
	}
	for n := 0; n < b.N; n++ {
		code, err, _ := ReadCommand(testDirectory, "test.txt", 0, 0)
		if code != returncodes.OK {
			log.Fatal("error reading test file:", err)
		}
	}
}

func BenchmarkStat(b *testing.B) {
	setupTestDirectory()
	code, err := CreatCommand(testDirectory, "test.txt", os.FileMode(0777), false)
	if code != returncodes.OK {
		log.Fatal("Error creating test file:", err)
	}
	for n := 0; n < b.N; n++ {
		code, err, _ := StatCommand(testDirectory, "test.txt")
		if code != returncodes.OK {
			log.Fatal("error stat test file:", err)
		}
	}
}

func BenchmarkAccess(b *testing.B) {
	setupTestDirectory()
	code, err := CreatCommand(testDirectory, "test.txt", os.FileMode(0777), false)
	if code != returncodes.OK {
		log.Fatal("error creating test file:", err)
	}
	for n := 0; n < b.N; n++ {
		code, err := AccessCommand(testDirectory, "test.txt", 0)
		if code != returncodes.OK {
			log.Fatal("error accessing test file:", err)
		}
	}
}

func BenchmarkTruncate(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		code, err := CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		if code != returncodes.OK {
			log.Fatal("error creating test file:", err)
		}
		code, err = TruncateCommand(testDirectory, "test.txt"+str, 3, false)
		if code != returncodes.OK {
			log.Fatal("error truncating test file:", err)
		}
	}
}

func BenchmarkUtimes(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		atime := time.Unix(100, 100)
		mtime := time.Unix(500, 250)
		code, err := CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		if code != returncodes.OK {
			log.Fatal("error creating test file:", err)
		}
		code, err = UtimesCommand(testDirectory, "test.txt"+str, &atime, &mtime, false)
		if code != returncodes.OK {
			log.Fatal("error changing test file time:", err)
		}
	}
}

func BenchmarkRmDir(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		code, err := MkdirCommand(testDirectory, "testDir"+str, os.FileMode(0777), false)
		if code != returncodes.OK {
			log.Fatal("error making benchdir:", err)
		}
		code, err = RmdirCommand(testDirectory, "testDir"+str, false)
		if code != returncodes.OK {
			log.Fatal("error removing benchdir:", err)
		}
	}
}

func BenchmarkReadDir(b *testing.B) {
	setupTestDirectory()
	code, err := MkdirCommand(testDirectory, "testDir", os.FileMode(0777), false)
	if code != returncodes.OK {
		log.Fatal("error making benchDir:", err)
	}
	code, err = CreatCommand(path.Join(testDirectory, "testDir"), "test.txt", os.FileMode(0777), false)
	if code != returncodes.OK {
		log.Fatal("error creating test file:", err)
	}
	for n := 0; n < b.N; n++ {
		code, err, _ = ReadDirCommand(testDirectory, "testDir")
		if code != returncodes.OK {
			log.Fatal("error reading benchDir:", err)
		}
	}
}

func BenchmarkLink(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		code, err := CreatCommand(testDirectory, "test.txt"+str, os.FileMode(0777), false)
		if code != returncodes.OK {
			log.Fatal("error creating testFile:", err)
		}
		code, err = LinkCommand(testDirectory, "test.txt"+str, "test2.txt"+str, false)
	}
}

func BenchmarkSymLink(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		code, err := SymlinkCommand(testDirectory, "testfolder/testlink", "testsymlink"+str, false)
		if code != returncodes.OK {
			log.Fatal("Symlink did not return OK. Actual:", code, " Error:", err)
		}
	}
}

func BenchmarkReadLink(b *testing.B) {
	setupTestDirectory()
	code, err := SymlinkCommand(testDirectory, "testfolder/testlink", "testsymlink", false)
	if code != returncodes.OK {
		log.Fatal("Symlink did not return OK. Actual:", code, " Error:", err)
	}
	for n := 0; n < b.N; n++ {
		code, err, _ = ReadlinkCommand(testDirectory, "testsymlink")
		if code != returncodes.OK {
			log.Fatalln("Readlink did not return OK. Actual:", code, " Error:", err)
		}
	}
}

func BenchmarkMkDir(b *testing.B) {
	setupTestDirectory()
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		code, err := MkdirCommand(testDirectory, "testDir"+str, os.FileMode(0777), false)
		if code != returncodes.OK {
			log.Fatal("error making benchdir:", err)
		}
	}
}
