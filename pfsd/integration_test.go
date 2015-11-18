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
	"testing"
)

func TestMain(m *testing.M) {
	os.Mkdir("/tmp/pfs", 0777)
	init := exec.Command("paranoid-cli", "init", "/tmp/pfs")
	init.Run()
	pnetserver.ParanoidDir = "/tmp/pfs"
	globals.Port = 10101
	lis, _ := net.Listen("tcp", ":10101")
	srv := grpc.NewServer()
	pb.RegisterParanoidNetworkServer(srv, &pnetserver.ParanoidServer{})
	go srv.Serve(lis)
	flag.Parse()
	exitCode := m.Run()
	os.RemoveAll("/tmp/pfs")
	os.Exit(exitCode)
}

func TestCreat(t *testing.T) {
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.CreatRequest{
		Path:        "creat.txt",
		Permissions: 0777,
	}
	_, err = client.Creat(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
}

func TestChmod(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", "/tmp/pfs", "chmod.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.ChmodRequest{
		Path: "chmod.txt",
		Mode: 0777,
	}
	_, err = client.Chmod(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
}

func TestLink(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", "/tmp/pfs", "link.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.LinkRequest{
		OldPath: "link.txt",
		NewPath: "linknew.txt",
	}
	_, err = client.Link(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
}

func TestPing(t *testing.T) {
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.PingRequest{
		Ip:   "0.0.0.0",
		Port: "0",
	}
	_, err = client.Ping(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
}

func TestRename(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", "/tmp/pfs", "rename.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.RenameRequest{
		OldPath: "rename.txt",
		NewPath: "renamenew.txt",
	}
	_, err = client.Rename(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
}

func TestTruncate(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", "/tmp/pfs", "truncate.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.TruncateRequest{
		Path:   "truncate.txt",
		Length: 0,
	}
	_, err = client.Truncate(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
}

func TestUnlink(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", "/tmp/pfs", "unlink.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.UnlinkRequest{
		Path: "unlink.txt",
	}
	_, err = client.Unlink(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
}

func TestUtimes(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", "/tmp/pfs", "utimes.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
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
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
}

func TestWrite(t *testing.T) {
	// Create file to run test on
	creat := exec.Command("pfsm", "creat", "/tmp/pfs", "write.txt", "0777")
	creat.Run()
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
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
		t.Errorf("Creat did not return OK. Actual: %v", err)
	}
	if resp.BytesWritten != 4 {
		t.Errorf("Incorrect bytes written. Expected: 4. Actual: %d", resp.BytesWritten)
	}
}
