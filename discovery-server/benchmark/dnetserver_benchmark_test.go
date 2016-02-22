package dnetservertest

import (
	. "github.com/cpssd/paranoid/discovery-server/dnetserver"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	syslog "log"
	"os"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	Log = logger.New("discoveryTest", "discoveryTest", "/dev/null")
	Log.SetLogLevel(logger.ERROR)
	os.Exit(m.Run())
}

func BenchmarkJoin(b *testing.B) {
	discovery := DiscoveryServer{}
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		joinRequest := pb.JoinRequest{
			Node: &pb.Node{CommonName: "TestNode" + str, Ip: "1.1.1." + str, Port: "1001", Uuid: "blahblah" + str},
			Pool: "TestPool",
		}
		_, err := discovery.Join(nil, &joinRequest)
		if err != nil {
			syslog.Fatalln("Error joining network : ", err)
		}
	}
}

func BenchmarkDisco(b *testing.B) {
	discovery := DiscoveryServer{}
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		joinRequest := pb.JoinRequest{
			Node: &pb.Node{CommonName: "TestNode" + str, Ip: "1.1.1.1" + str, Port: "1001", Uuid: "blahblah"},
			Pool: "TestPool",
		}
		discovery.Join(nil, &joinRequest)
		disconnect := pb.DisconnectRequest{
			Node: &pb.Node{CommonName: "TestNode" + str, Ip: "1.1.1.1" + str, Port: "1001", Uuid: "blahblah"},
		}
		_, err := discovery.Disconnect(nil, &disconnect)
		if err != nil {
			syslog.Fatalln("Error disconnecting node 2:", err)
		}
	}
}
