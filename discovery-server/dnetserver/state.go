package dnetserver

import (
	"encoding/json"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"io/ioutil"
	"os"
	"time"
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

	var jsonNodes []jsonNode
	err = json.Unmarshal(fileData, &jsonNodes)
	if err != nil {
		Log.Fatal("Failed to un-marshal state file:", err)
	}

	newNodes := make([]Node, len(jsonNodes))
	for i, val := range jsonNodes {
		newNodes[i] = regularNodeFromJsonNode(val)
	}

	Nodes = newNodes
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

// getJsonReadyNodeList returns a list of jsonNodes based on the current Nodes
func getJsonReadyNodeList() []jsonNode {
	jsonNodes := make([]jsonNode, len(Nodes))
	for i, val := range Nodes {
		jsonNodes[i] = jsonNodeFromRegularNode(val)
	}
	return jsonNodes
}

// jsonNodeFromRegularNode returns a jsonNode based on the given Node
func jsonNodeFromRegularNode(n Node) jsonNode {
	return jsonNode{
		Pool:        n.Pool,
		ExpiryTime:  n.ExpiryTime.Unix(),
		IP:          n.Data.Ip,
		Port:        n.Data.Port,
		Common_name: n.Data.CommonName,
		UUID:        n.Data.Uuid,
	}
}

// regularNodeFromJsonNode returns a Node based on the given jsonNode
func regularNodeFromJsonNode(n jsonNode) Node {
	return Node{
		Pool:       n.Pool,
		ExpiryTime: time.Unix(n.ExpiryTime, 0),
		Data: pb.Node{
			Ip:         n.IP,
			Port:       n.Port,
			CommonName: n.Common_name,
			Uuid:       n.UUID,
		},
	}
}
