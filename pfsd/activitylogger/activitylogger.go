package activitylogger

import (
	"github.com/cpssd/paranoid/logger"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

var (
	logDir       string
	currentIndex int
	indexLock    sync.Mutex
	pLog         *logger.ParanoidLogger
)

// Init initialises the logger
func Init(paranoidDirectory string) {
	logDir = path.Join(paranoidDirectory, "meta", "activity_logs")
	pLog = logger.New("Activity Logger", "pfsd", path.Join(paranoidDirectory, "meta", "logs"))
	fileDescriptors, err := ioutil.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			setup(logDir)
			return
		} else if os.IsPermission(err) {
			pLog.Fatal("Activity logger does not have permissions for: ", logDir)
		} else {
			pLog.Fatal(err)
		}
	}
	currentIndex = len(fileDescriptors) + 1000000
}

// setup is called when the log directory does not exist
func setup(logDirectory string) {
	err := os.MkdirAll(logDirectory, 0777)
	if err != nil {
		pLog.Fatal("failed to create log directory")
	}
	currentIndex = 1000000
}

// LastEntryIndex returns the index of the last log entry
func LastEntryIndex() int {
	indexLock.Lock()
	defer indexLock.Unlock()
	if currentIndex == 1000000 {
		return -1
	}
	return currentIndex - 1
}
