// +build !integration

package dnetserver

import (
	"encoding/json"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	Log = logger.New("discoveryTest", "discoveryTest", "/dev/null")
	StateFilePath = path.Join(os.TempDir(), "server_state.json")
	os.Exit(m.Run())
}

func TestStateSave(t *testing.T) {
	discovery := DiscoveryServer{}

	//Join node1
	joinRequest := pb.JoinRequest{
		Node: &pb.Node{CommonName: "TestNode1", Ip: "1.1.1.1", Port: "1001", Uuid: "blahblah1"},
		Pool: "TestPool",
	}
	joinResponse, err := discovery.Join(nil, &joinRequest)
	if err != nil {
		t.Error("Error joining network : ", err)
	}
	if len(joinResponse.Nodes) != 0 {
		t.Error("Incorrect nodes returned :", joinResponse.Nodes)
	}

	stateFileData, err := ioutil.ReadFile(StateFilePath)
	if err != nil {
		t.Error("Failed to read state file: ", err)
	}

	var jsonNodes []jsonNode
	err = json.Unmarshal(stateFileData, &jsonNodes)
	if err != nil {
		Log.Fatal("Failed to un-marshal state file:", err)
	}

	if len(jsonNodes) != 1 {
		t.Error("wrong number of nodes in state file:", len(jsonNodes))
	}
	if jsonNodes[0].UUID != "blahblah1" || jsonNodes[0].Pool != "TestPool" {
		t.Error("Node in state file is wrong: ", jsonNodes[0])
	}
}

func TestStateLoad(t *testing.T) {
	LoadState()

	if len(Nodes) != 1 {
		t.Error("Wrong number of nodes loaded from state file")
	}

	if Nodes[0].Pool != "TestPool" || Nodes[0].Data.Uuid != "blahblah1" {
		t.Error("loaded node is wrong:", Nodes[0])
	}

	Nodes = nil
}

func TestDiscoveryNetwork(t *testing.T) {
	discovery := DiscoveryServer{}

	//Join node1
	joinRequest := pb.JoinRequest{
		Node: &pb.Node{CommonName: "TestNode1", Ip: "1.1.1.1", Port: "1001", Uuid: "blahblah1"},
		Pool: "TestPool",
	}
	joinResponse, err := discovery.Join(nil, &joinRequest)
	if err != nil {
		t.Error("Error joining network : ", err)
	}
	if len(joinResponse.Nodes) != 0 {
		t.Error("Incorrect nodes returned :", joinResponse.Nodes)
	}

	//Join node2
	joinRequest = pb.JoinRequest{
		Node: &pb.Node{CommonName: "TestNode2", Ip: "1.1.1.2", Port: "1001", Uuid: "blahblah2"},
		Pool: "TestPool",
	}
	joinResponse, err = discovery.Join(nil, &joinRequest)
	if err != nil {
		t.Error("Error joining network : ", err)
	}
	if len(joinResponse.Nodes) != 1 {
		t.Error("Incorrect nodes returned :", joinResponse.Nodes)
	}
	if joinResponse.Nodes[0].CommonName != "TestNode1" || joinResponse.Nodes[0].Ip != "1.1.1.1" ||
		joinResponse.Nodes[0].Port != "1001" {
		t.Error("Incorrect node information returned :", joinResponse.Nodes[0])
	}

	//Join node3
	joinRequest = pb.JoinRequest{
		Node: &pb.Node{CommonName: "TestNode3", Ip: "1.1.1.1", Port: "1002", Uuid: "blahblah3"},
		Pool: "TestPool",
	}
	joinResponse, err = discovery.Join(nil, &joinRequest)
	if err != nil {
		t.Error("Error joining network : ", err)
	}
	if len(joinResponse.Nodes) != 2 {
		t.Error("Incorrect nodes returned :", joinResponse.Nodes)
	}
	expectedNode1 := pb.Node{CommonName: "TestNode1", Ip: "1.1.1.1", Port: "1001", Uuid: "blahblah1"}
	expectedNode2 := pb.Node{CommonName: "TestNode2", Ip: "1.1.1.2", Port: "1001", Uuid: "blahblah2"}
	if (*joinResponse.Nodes[0] != expectedNode1 || *joinResponse.Nodes[1] != expectedNode2) &&
		(*joinResponse.Nodes[0] != expectedNode2 || *joinResponse.Nodes[1] != expectedNode1) {
		t.Error("Incorrect nodes returned :", joinResponse.Nodes)
	}

	//Disconnect node2
	disconnectRequest := pb.DisconnectRequest{
		Node: &pb.Node{CommonName: "TestNode2", Ip: "1.1.1.2", Port: "1001", Uuid: "blahblah2"},
	}
	_, err = discovery.Disconnect(nil, &disconnectRequest)
	if err != nil {
		t.Error("Error disconnecting node 2:", err)
	}

	//Join node2 (again)
	joinRequest = pb.JoinRequest{
		Node: &pb.Node{CommonName: "TestNode2", Ip: "1.1.1.2", Port: "1001", Uuid: "blahblah2"},
		Pool: "TestPool",
	}
	joinResponse, err = discovery.Join(nil, &joinRequest)
	if err != nil {
		t.Error("Error joining network : ", err)
	}
	if len(joinResponse.Nodes) != 2 {
		t.Error("Incorrect nodes returned :", joinResponse.Nodes)
	}
	expectedNode1 = pb.Node{CommonName: "TestNode1", Ip: "1.1.1.1", Port: "1001", Uuid: "blahblah1"}
	expectedNode2 = pb.Node{CommonName: "TestNode3", Ip: "1.1.1.1", Port: "1002", Uuid: "blahblah3"}
	if (*joinResponse.Nodes[0] != expectedNode1 || *joinResponse.Nodes[1] != expectedNode2) &&
		(*joinResponse.Nodes[0] != expectedNode2 || *joinResponse.Nodes[1] != expectedNode1) {
		t.Error("Incorrect nodes returned :", joinResponse.Nodes)
	}

	//Join node4
	joinRequest = pb.JoinRequest{
		Node: &pb.Node{CommonName: "TestNode4", Ip: "1.1.1.3", Port: "1001", Uuid: "blahblah4"},
		Pool: "TestPool",
	}
	joinResponse, err = discovery.Join(nil, &joinRequest)
	if err != nil {
		t.Error("Error joining network : ", err)
	}
	if len(joinResponse.Nodes) != 3 {
		t.Error("Incorrect nodes returned :", joinResponse.Nodes)
	}
}
