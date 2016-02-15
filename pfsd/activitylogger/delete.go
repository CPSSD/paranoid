package activitylogger

import (
	"errors"
	"log"
	"os"
	"path"
	"strconv"
)

// DeleteEntry deletes an entry in the logs per index and all logs after it
func DeleteEntry(index int) error {
	pause()
	defer resume()
	if index < 0 || index >= currentIndex {
		return errors.New("Index out of bounds")
	}

	for i := currentIndex - 1; i >= index; i-- {
		err := os.Remove(path.Join(logDir, strconv.Itoa(i)))
		if err != nil {
			log.Fatalln("Activity logger: failed to delete log of index:", index, "with error:", err)
		}
		currentIndex--
	}

	return nil
}
