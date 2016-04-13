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

	newStateFilePath := path.Join(TempDirectoryPath, pool)
	stateFilePath := path.Join(StateDirectoryPath, pool)

	err = ioutil.WriteFile(newStateFilePath, stateData, 0600)
	if err != nil {
		Log.Fatal("Failed to write state data to state file:", err)
	}

	err = os.Rename(newStateFilePath, stateFilePath)
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
