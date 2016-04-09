package dnetserver

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func saveState(pool string) {
	stateData, err := json.Marshal(Pools[pool].Info)
	if err != nil {
		Log.Fatal("Couldn't marshal stateData:", err)
	}

	file := prepareStateFile(pool)
	defer file.Close()
	_, err = file.Write(stateData)
	if err != nil {
		Log.Fatal("Failed to write state data to state file:", err)
	}
}

// LoadState loads the state from the statefiles in the state directory
func LoadState() {
	_, err := os.Stat(StateDirectoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			Log.Info("Tried loading state from state directory but it's non-existant")
			return
		} else {
			Log.Fatal("Couldn't stat state directory:", err)
		}
	}

	files, err := ioutil.ReadDir(StateDirectoryPath)
	if err != nil {
		Log.Fatal("Couldn't read state directory:", err)
	}

	for i := 0; i < len(files); i++ {
		stateFilePath := path.Join(StateDirectoryPath, files[i].Name())
		fileData, err := ioutil.ReadFile(stateFilePath)
		if err != nil {
			Log.Fatalf("Couldn't read state file %s: %s", stateFilePath, err)
		}

		Pools[files[i].Name()] = &Pool{
			Info: PoolInfo{},
		}

		err = json.Unmarshal(fileData, &Pools[files[i].Name()].Info)
		if err != nil {
			Log.Fatalf("Failed to un-marshal state file %s: %s", stateFilePath, err)
		}
	}
}

// prepareStateFile prepares the statefile for a state update and returns the file
// ready to be writen to
func prepareStateFile(pool string) *os.File {
	stateFilePath := path.Join(StateDirectoryPath, pool)
	_, err := os.Stat(stateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			Log.Info("Creating state file: ", stateFilePath)
			_, err1 := os.Create(stateFilePath)
			if err1 != nil {
				Log.Fatal("Could not create state file:", err1)
			}
		} else {
			Log.Fatal("Couldnt stat stateFile:", err)
		}
	}

	err = os.Truncate(stateFilePath, 0)
	if err != nil {
		Log.Fatal("Couldnt truncate state file:", err)
	}

	file, err := os.OpenFile(stateFilePath, os.O_WRONLY, 0600)
	if err != nil {
		Log.Fatal("Couldnt open stateFile:", err)
	}

	return file
}
