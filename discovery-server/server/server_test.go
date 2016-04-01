// +build !integration

package server_test

import (
	"github.com/cpssd/paranoid/discovery-server/dnetserver"
	. "github.com/cpssd/paranoid/discovery-server/server"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	fileserver "github.com/cpssd/paranoid/proto/fileserver"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	FileMap = make(map[string]*FileCache)
	os.Exit(m.Run())
}

func testFileShare(t *testing.T) {
	discovery := dnetserver.DiscoveryServer{}
	server := FileserverServer{}

	joinRequest := pb.JoinRequest{
		Node: &pb.Node{CommonName: "TestNode1", Ip: "1.1.1.1", Port: "1001", Uuid: "blahblah1"},
		Pool: "TestPool",
	}
	_, err := discovery.Join(nil, &joinRequest)
	if err != nil {
		t.Error("Error joining network : ", err)
	}

	request := fileserver.ServeRequest{
		Uuid:     "blahblah1",
		FilePath: "asdf.txt",
		FileData: []byte("This is a Test"),
	}

	response, err := server.ServeFile(nil, &request)
	if err != nil {
		t.Error("Error adding File")
	}
	_, ok := FileMap[response.ServeResponse]
	if !ok {
		t.Error("Failure Caching File")
	}
}
