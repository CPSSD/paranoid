package activitylogger

import (
	"github.com/cpssd/paranoid/logger"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

// ActivityLogger is the structure through which logging functinality can be accessed
type ActivityLogger struct {
	logDir       string
	currentIndex uint64
	indexLock    sync.Mutex
	pLog         *logger.ParanoidLogger
}

// New returns an initiated instance of ActivityLogger
func New(paranoidDirectory string) *ActivityLogger {
	al := &ActivityLogger{
		logDir: path.Join(paranoidDirectory, "meta", "activity_logs"),
		pLog:   logger.New("Activity Logger", "pfsd", path.Join(paranoidDirectory, "meta", "logs")),
	}
	fileDescriptors, err := ioutil.ReadDir(al.logDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(al.logDir, 0777)
			if err != nil {
				al.pLog.Fatal("failed to create log directory")
			}
		} else if os.IsPermission(err) {
			al.pLog.Fatal("Activity logger does not have permissions for: ", al.logDir)
		} else {
			al.pLog.Fatal(err)
		}
	}
	al.currentIndex = uint64(len(fileDescriptors) + 1)
	return al
}

// LastEntryIndex returns the index of the last log entry
func (al *ActivityLogger) LastEntryIndex() uint64 {
	al.indexLock.Lock()
	defer al.indexLock.Unlock()
	return al.currentIndex - 1
}

// f12ci converts a fileIndex to a convenient index
func fi2ci(i uint64) uint64 {
	return i - 1000000
}

// ci2fi converts a convenient to a fileIndex
func ci2fi(i uint64) uint64 {
	return i + 1000000
}
