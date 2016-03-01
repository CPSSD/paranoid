package raftlog

import (
	"errors"
	"os"
	"path"
	"strconv"
)

// DiscardLogEntries an entry in the logs per index and all logs after it
func (rl *RaftLog) DiscardLogEntries(startIndex uint64) error {
	rl.indexLock.Lock()
	defer rl.indexLock.Unlock()

	if startIndex < 1 || startIndex >= rl.currentIndex {
		return errors.New("Index out of bounds")
	}

	for i := rl.currentIndex - 1; i >= startIndex; i-- {
		err := os.Remove(path.Join(rl.logDir, strconv.FormatUint(storageIndexToFileIndex(i), 10)))
		if err != nil {
			rl.pLog.Fatal("Activity logger: failed to delete log of index:",
				i, "with error:", err)
		}
		rl.currentIndex--
	}

	if rl.currentIndex > 1 {
		logEntry, err := rl.GetLogEntry(rl.currentIndex - 1)
		if err != nil {
			rl.pLog.Fatal("Error deleting log entries:", err)
		}
		rl.mostRecentTerm = logEntry.Term
	} else {
		rl.mostRecentTerm = 0
	}

	return nil
}
