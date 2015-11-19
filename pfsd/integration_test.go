// +build integration

package main

import (
	"flag"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net"
	"os"
	"os/exec"
	"path"
	"testing"
)

var tmpdir = path.Join(os.TempDir(), "pfs")

func TestMain(m *testing.M) {
	// Make sure we start with an empty directory
	os.RemoveAll(tmpdir)
	os.Mkdir(tmpdir, 0777)
	init := exec.Command("pfsm", "init", tmpdir)
	init.Run()
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

func TestCreat(t *testing.T) {
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.CreatRequest{
		Path:        "creat.txt",
		Permissions: 0777,
	}
	_, err = client.Creat(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual:", err)
	}
}

func TestChmod(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", tmpdir, "chmod.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.ChmodRequest{
		Path: "chmod.txt",
		Mode: 0777,
	}
	_, err = client.Chmod(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Chmod did not return OK. Actual:", err)
	}
}

func TestLink(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", tmpdir, "link.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.LinkRequest{
		OldPath: "link.txt",
		NewPath: "linknew.txt",
	}
	_, err = client.Link(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Link did not return OK. Actual:", err)
	}
}

func TestPing(t *testing.T) {
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.PingRequest{
		Ip:   "0.0.0.0",
		Port: "0",
	}
	_, err = client.Ping(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Ping did not return OK. Actual:", err)
	}
}

func TestRename(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", tmpdir, "rename.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.RenameRequest{
		OldPath: "rename.txt",
		NewPath: "renamenew.txt",
	}
	_, err = client.Rename(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Rename did not return OK. Actual:", err)
	}
}

func TestTruncate(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", tmpdir, "truncate.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.TruncateRequest{
		Path:   "truncate.txt",
		Length: 0,
	}
	_, err = client.Truncate(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Error("Truncate did not return OK. Actual:", err)
	}
}

func TestUnlink(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", tmpdir, "unlink.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatal("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.UnlinkRequest{
		Path: "unlink.txt",
	}
	_, err = client.Unlink(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Error("Unlink did not return OK. Actual:", err)
	}
}

func TestUtimes(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", tmpdir, "utimes.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatal("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.UtimesRequest{
		Path:               "utimes.txt",
		AccessSeconds:      1,
		AccessMicroseconds: 1,
		ModifySeconds:      1,
		ModifyMicroseconds: 1,
	}
	_, err = client.Utimes(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Error("Utimes did not return OK. Actual:", err)
	}
}

func TestWrite(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", tmpdir, "write.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatal("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.WriteRequest{
		Path:   "write.txt",
		Data:   []byte("YmxhaA=="), // "blah"
		Offset: 0,
		Length: 4,
	}
	resp, err := client.Write(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Error("Write did not return OK. Actual:", err)
	}
	if resp.BytesWritten != 4 {
		t.Error("Incorrect bytes written. Expected: 4. Actual:", resp.BytesWritten)
	}
}
