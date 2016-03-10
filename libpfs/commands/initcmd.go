package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//makeDir creates a new directory with permissions 0777 with the name newDir in parentDir.
func makeDir(parentDir, newDir string) (string, error) {
	newDirPath := path.Join(parentDir, newDir)
	err := os.Mkdir(newDirPath, 0700)
	if err != nil {
		return "", err
	}
	return newDirPath, nil
}

//checkEmpty checks if a given directory has any children.
func checkEmpty(directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("error reading directory", err)
	}
	if len(files) > 0 {
		return errors.New("init : directory must be empty")
	}
	return nil
}

//InitCommand creates the pvd directory sturucture
//It also gets a UUID and stores it in the meta directory.
func InitCommand(directory string) (returnCode int, returnError error) {
	err := checkEmpty(directory)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}
	Log.Verbose("init : creating new paranoid file system in " + directory)

	_, err = makeDir(directory, "names")
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	_, err = makeDir(directory, "inodes")
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	metaDir, err := makeDir(directory, "meta")
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	_, err = makeDir(metaDir, "logs")
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	_, err = makeDir(metaDir, "raft")
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	_, err = makeDir(directory, "contents")
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	uuid, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading uuid:", err)
	}

	uuidString := strings.TrimSpace(string(uuid))
	Log.Verbose("init uuid : " + uuidString)

	err = ioutil.WriteFile(path.Join(metaDir, "uuid"), []byte(uuidString), 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing uuid file:", err)
	}

	_, err = os.Create(path.Join(metaDir, "lock"))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating lock file:", err)
	}
	return returncodes.OK, nil
}
