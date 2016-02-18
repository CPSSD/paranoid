package activitylogger

import (
	"errors"
	"os"
	"path"
	"strconv"
)

// DeleteEntry deletes an entry in the logs per index and all logs after it
func (al *ActivityLogger) DeleteEntry(index uint64) error {
	al.indexLock.Lock()
	defer al.indexLock.Unlock()

	if index < 1 || index >= al.currentIndex {
		return errors.New("Index out of bounds")
	}

	for i := al.currentIndex - 1; i >= index; i-- {
		err := os.Remove(path.Join(al.logDir, strconv.FormatUint(ci2fi(i), 10)))
		if err != nil {
			al.pLog.Fatal("Activity logger: failed to delete log of index:",
				i, "with error:", err)
		}
		al.currentIndex--
	}

	return nil
}
