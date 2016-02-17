// +build !integration benchmark

package main

import (
	"flag"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"os"
	"path"
	"strconv"
	"testing"
)

var tmpdir = path.Join(os.TempDir(), "pfs")

func TestMain(m *testing.M) {
	pnetserver.Log = logger.New("pnetserver", "pfsd", os.DevNull)
	pnetserver.Log.SetLogLevel(logger.ERROR)
	os.RemoveAll(tmpdir)
	os.Mkdir(tmpdir, 0777)
	commands.Log = logger.New("pfsdBench", "pfsdBench", os.DevNull)
	commands.Log.SetLogLevel(logger.ERROR)
	commands.InitCommand(tmpdir)
	pnetserver.ParanoidDir = tmpdir
	globals.Port = 10101
	lis, _ := net.Listen("tcp", ":10101")
	srv := grpc.NewServer()
	pb.RegisterParanoidNetworkServer(srv, &pnetserver.ParanoidServer{})
	go srv.Serve(lis)
	flag.Parse()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func createServer() *grpc.ClientConn {
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	return conn
}

func BenchmarkPing(b *testing.B) {
	conn := createServer()
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		req := pb.PingRequest{
			Ip:   "0.0.0.0",
			Port: "0",
		}
		client.Ping(context.Background(), &req)
	}
}

func BenchmarkCreat(b *testing.B) {
	conn := createServer()
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n + 1000)
		rand := strconv.Itoa(rand.Intn(100000))
		req := pb.CreatRequest{
			Path:        "test.txt" + str + rand,
			Permissions: 0777,
		}
		client.Creat(context.Background(), &req)
	}
}

func BenchmarkWrite(b *testing.B) {
	conn := createServer()
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
		client.Write(context.Background(), &req)
	}
}

func BenchmarkChmod(b *testing.B) {
	conn := createServer()
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.CreatCommand(tmpdir, "chmodBench.txt"+str, os.FileMode(0744), false)
		req := pb.ChmodRequest{
			Path: "chmodBench.txt" + str,
			Mode: 0777,
		}
		client.Chmod(context.Background(), &req)
	}
}

func BenchmarkUtimes(b *testing.B) {
	conn := createServer()
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
		client.Utimes(context.Background(), &req)
	}
}

func BenchmarkTruncate(b *testing.B) {
	conn := createServer()
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.CreatCommand(tmpdir, "truncateBench.txt"+str, os.FileMode(0777), false)
		req := pb.TruncateRequest{
			Path:   "truncateBench.txt" + str,
			Length: 0,
		}
		client.Truncate(context.Background(), &req)
	}
}

func BenchmarkRename(b *testing.B) {
	conn := createServer()
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
		client.Rename(context.Background(), &req)
	}
}

func BenchmarkLink(b *testing.B) {
	conn := createServer()
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
		client.Link(context.Background(), &req)
	}
}

func BenchmarkMkdir(b *testing.B) {
	conn := createServer()
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		req := pb.MkdirRequest{
			Directory: "somedir" + str,
			Mode:      0777,
		}
		client.Mkdir(context.Background(), &req)
	}
}

func BenchmarkRmdir(b *testing.B) {
	conn := createServer()
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.MkdirCommand(tmpdir, "rmDirBench"+str, os.FileMode(0777), false)
		req := pb.RmdirRequest{
			Directory: "rmDirBench" + str,
		}
		client.Rmdir(context.Background(), &req)
	}
}
