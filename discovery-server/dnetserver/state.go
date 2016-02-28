package dnetserver

import (
	"encoding/json"
	"os"
)

// a json marshalable structure for a node
type jsonNode struct {
	Pool        string `json:"pool"`
	ExpiryTime  int64  `json:"expiryTime"`
	IP          string `json:"ip"`
	Port        string `json:"port"`
	Common_name string `json:"common_name"`
	UUID        string `jsoin:"uuid"`
}

// saveState saves the current state of the discovery server to a file in it's
// meta directory. the state includes all pools and nodes
func saveState() {
	currentNodes := getJsonReadyNodeList()
	stateData, err := json.Marshal(currentNodes)
	if err != nil {
		Log.Fatal("Couldnt marshal stateData:", err)
	}

	file := prepareStateFile()
	defer file.Close()
	_, err = file.Write(stateData)
	if err != nil {
		Log.Fatal("Failed to write state date to state file:", err)
	}
}

// prepareStateFile prepares the statefile for a state update and returns the file
// ready to be writen to
func prepareStateFile() *os.File {
	_, err := os.Stat(StateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			Log.Info("Creating state file: ", StateFilePath)
			file, err1 := os.Create(StateFilePath)
			if err1 != nil {
				Log.Fatal("Could not create state file:", err1)
			}
			return file
		} else {
			Log.Fatal("Couldnt stat stateFile:", err)
		}
	}

	err = os.Truncate(StateFilePath, 0)
	if err != nil {
		Log.Fatal("Couldnt truncate state file:", err)
	}

	file, err := os.Open(StateFilePath)
	if err != nil {
		Log.Fatal("Couldnt open stateFile:", err)
	}

	return file
}

// getJsonReadyNodeList returns a list of jsonNodes based on the current Nodes
func getJsonReadyNodeList() []jsonNode {
	jsonNodes := make([]jsonNode, len(Nodes))
	for i, val := range Nodes {
		jsonNodes[i] = JsonNodeFromRegularNode(val)
	}
	return jsonNodes
}

// JsonNodeFromRegularNode returns a jsonNode based on the given Node
func JsonNodeFromRegularNode(n Node) jsonNode {
	return jsonNode{
		Pool:        n.Pool,
		ExpiryTime:  n.ExpiryTime.Unix(),
		IP:          n.Data.Ip,
		Port:        n.Data.Port,
		Common_name: n.Data.CommonName,
		UUID:        n.Data.Uuid,
	}
}
