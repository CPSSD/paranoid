package activitylogger

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

var (
	logDir        string
	currentIndex  int
	appendLogChan chan logEntry
	killChan      chan bool
	pauseChan     chan bool
	resumeChan    chan bool
	paused        bool
)

const (
	typeChmod uint8 = iota
	typeCreat
	typeLink
	typeMkdir
	typeRename
	typeRmdir
	typeSymLink
	typeTruncate
	typeUnlink
	typeUtimes
	typeWrite
)

// logEntry is an abstaction of a log entry to be passed through the appendLog channel
type logEntry struct {
	EntryType uint8
	Params    []interface{}
}

func newLogEntry(typ uint8, params ...interface{}) logEntry {
	return logEntry{typ, params}
}

// Init initialises the logger
func Init(paranoidDirectory string) {
	logDir = path.Join(paranoidDirectory, "meta", "activity_logs")
	fileDescriptors, err := ioutil.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			err1 := os.Mkdir(logDir, 0777)
			if err1 != nil {
				log.Fatalln(err1)
			}
		} else if os.IsPermission(err) {
			log.Fatalln("Activity logger does not have permissions for: ", logDir)
		} else {
			log.Fatalln(err)
		}
	}

	if len(fileDescriptors) == 0 {
		currentIndex = 0
	} else {
		currentIndex, err = strconv.Atoi(fileDescriptors[len(fileDescriptors)-1].Name())
		if err != nil {
			log.Fatalln(err)
		}
		currentIndex++
	}

	appendLogChan = make(chan logEntry, 200)
	killChan = make(chan bool, 1)
	pauseChan = make(chan bool, 1)
	resumeChan = make(chan bool, 1)
	go listenAppendLog()
}

// LastLogIndex returns the index
func LastLogIndex() int {
	if currentIndex == 0 {
		return -1
	}
	return currentIndex
}

func Get(index int) {

}
