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
	"math/rand"
	"net"
	"os"
	"path"
	"strconv"
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

func BenchmarkPFSDCreat(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n + 1000)
		rand := strconv.Itoa(rand.Intn(100000))
		req := pb.CreatRequest{
			Path:        "tessyslog.txt" + str + rand,
			Permissions: 0777,
		}
		_, err = client.Creat(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Creat did not return OK. Actual:", err)
		}
	}
}

func BenchmarkPFSDWrite(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.CreatCommand(tmpdir, "writeBench.txt"+str, os.FileMode(0777), false)
		req := pb.WriteRequest{
			Path:   "writeBench.txt" + str,
			Data:   []byte("Hello World"),
			Offset: 0,
			Length: 10,
		}
		_, err = client.Write(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Write did not return OK. Actual:", err)
		}
	}
}

func BenchmarkPFSDChmod(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.CreatCommand(tmpdir, "chmodBench.txt"+str, os.FileMode(0744), false)
		req := pb.ChmodRequest{
			Path: "chmodBench.txt" + str,
			Mode: 0777,
		}
		_, err = client.Chmod(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Chmod did not return OK. Actual:", err)
		}
	}
}

func BenchmarkPFSDUtimes(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.CreatCommand(tmpdir, "utimeBench.txt"+str, os.FileMode(0777), false)
		req := pb.UtimesRequest{
			Path:              "utimeBench.txt" + str,
			AccessSeconds:     1,
			AccessNanoseconds: 1,
			ModifySeconds:     1,
			ModifyNanoseconds: 1,
		}
		_, err = client.Utimes(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Utimes did not return OK. Actual:", err)
		}
	}
}

func BenchmarkPFSDTruncate(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.CreatCommand(tmpdir, "truncateBench.txt"+str, os.FileMode(0777), false)
		req := pb.TruncateRequest{
			Path:   "truncateBench.txt" + str,
			Length: 0,
		}
		_, err = client.Truncate(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Truncate did not return OK. Actual:", err)
		}
	}
}

func BenchmarkPFSDRename(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		rand := strconv.Itoa(rand.Intn(100000))
		commands.CreatCommand(tmpdir, "garbledNameBench.txt"+str+rand, os.FileMode(0777), false)
		req := pb.RenameRequest{
			OldPath: "garbledNameBench.txt" + str + rand,
			NewPath: "renameBench.txt" + str + rand,
		}
		_, err = client.Rename(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("rename did not return OK. Actual:", err)
		}
	}
}

func BenchmarkPFSDLink(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		rand := strconv.Itoa(rand.Intn(100000))
		commands.CreatCommand(tmpdir, "linkBench.txt"+str+rand, os.FileMode(0777), false)
		req := pb.LinkRequest{
			OldPath: "linkBench.txt" + str + rand,
			NewPath: "newlinkBench.txt" + str + rand,
		}
		_, err = client.Link(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Link did not return OK. Actual:", err)
		}
	}
}

func BenchmarkPFSDMkdir(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		rand := strconv.Itoa(rand.Intn(100000))
		str := strconv.Itoa(n)
		req := pb.MkdirRequest{
			Directory: "somedir" + str + rand,
			Mode:      0777,
		}
		_, err = client.Mkdir(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Mkdir did not return OK. Actual:", err)
		}
	}
}

func BenchmarkPFSDRmdir(b *testing.B) {
	conn, err := grpc.Dial("localhost:10102", grpc.WithInsecure())
	if err != nil {
		syslog.Fatal("Error Dialing server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.MkdirCommand(tmpdir, "rmDirBench"+str, os.FileMode(0777), false)
		req := pb.RmdirRequest{
			Directory: "rmDirBench" + str,
		}
		_, err = client.Rmdir(context.Background(), &req)
		if grpc.Code(err) != codes.OK {
			syslog.Fatal("Rmdir did not return OK. Actual:", err)
		}
	}
}
