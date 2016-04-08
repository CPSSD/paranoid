package dnetserver

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type PersistentState struct {
	Nodes []Node               `json:"nodes"`
	Pools map[string]*PoolInfo `json:"pools"`
}

// saveState saves the current state of the discovery server to a file in it's
// meta directory. the state includes all pools and nodes
func saveState() {
	stateData, err := json.Marshal(&PersistentState{
		Nodes: Nodes,
		Pools: Pools,
	})

	if err != nil {
		Log.Fatal("Couldn't marshal stateData:", err)
	}

	file := prepareStateFile()
	defer file.Close()
	_, err = file.Write(stateData)
	if err != nil {
		Log.Fatal("Failed to write state data to state file:", err)
	}
}

// LoadState loads the Nodes from the stateFile
func LoadState() {
	_, err := os.Stat(StateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			Log.Info("Tried loading state from state file but it's non-existant")
			return
		} else {
			Log.Fatal("Couldn't stat statefile:", err)
		}
	}

	fileData, err := ioutil.ReadFile(StateFilePath)

	perState := &PersistentState{}

	err = json.Unmarshal(fileData, perState)
	if err != nil {
		Log.Fatal("Failed to un-marshal state file:", err)
	}

	Nodes = perState.Nodes
	Pools = perState.Pools
}

// prepareStateFile prepares the statefile for a state update and returns the file
// ready to be writen to
func prepareStateFile() *os.File {
	_, err := os.Stat(StateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			Log.Info("Creating state file: ", StateFilePath)
			_, err1 := os.Create(StateFilePath)
			if err1 != nil {
				Log.Fatal("Could not create state file:", err1)
			}
		} else {
			Log.Fatal("Couldnt stat stateFile:", err)
		}
	}

	err = os.Truncate(StateFilePath, 0)
	if err != nil {
		Log.Fatal("Couldnt truncate state file:", err)
	}

	file, err := os.OpenFile(StateFilePath, os.O_WRONLY, 0600)
	if err != nil {
		Log.Fatal("Couldnt open stateFile:", err)
	}

	return file
}
