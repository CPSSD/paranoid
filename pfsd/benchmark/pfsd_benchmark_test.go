package benchmarkpfsd

import (
	"flag"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	syslog "log"
	"net"
	"os"
	"path"
	"testing"
)

var tmpdir = path.Join(os.TempDir(), "pfs")

func TestMain(m *testing.M) {
	// Make sure we start with an empty directory
	pnetserver.Log = logger.New("pnetserver", "pfsd", os.DevNull)
	pnetserver.Log.SetLogLevel(logger.ERROR)
	os.RemoveAll(tmpdir)
	os.Mkdir(tmpdir, 0777)
	commands.Log = logger.New("pfsdintegration", "pfsdintegration", os.DevNull)
	commands.InitCommand(tmpdir)
	commands.Log.SetLogLevel(logger.ERROR)
	pnetserver.ParanoidDir = tmpdir
	globals.Port = 10102
	lis, err := net.Listen("tcp", ":10102")
	if err != nil {
		syslog.Fatal("Error Creating PFSD server:", err)
	}
	srv := grpc.NewServer()
	pb.RegisterParanoidNetworkServer(srv, &pnetserver.ParanoidServer{})
	go srv.Serve(lis)
	flag.Parse()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func BenchmarkPFSDPing(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		req := pb.PingRequest{
			Ip:   "0.0.0.0",
			Port: "0",
		}
		_, err = client.Ping(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Ping did not return OK. Actual:", err)
		}
	}
}
