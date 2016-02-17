// +build integration

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
	"google.golang.org/grpc/codes"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"testing"
	"time"
)

var tmpdir = path.Join(os.TempDir(), "pfs")

func TestMain(m *testing.M) {
	// Make sure we start with an empty directory
	pnetserver.Log = logger.New("pnetserver", "pfsd", os.DevNull)
	os.RemoveAll(tmpdir)
	os.Mkdir(tmpdir, 0777)
	commands.Log = logger.New("pfsdintegration", "pfsdintegration", os.DevNull)
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

func TestCreat(t *testing.T) {
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %s", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.CreatRequest{
		Path:        "creat.txt",
		Permissions: 0777,
	}
	_, err = client.Creat(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Creat did not return OK. Actual: %s", err)
	}
}

func TestChmod(t *testing.T) {
	// Create file to run test on
	_, err := commands.CreatCommand(tmpdir, "chmod.txt", os.FileMode(0777), false)
	if err != nil {
		t.Fatalf("Could not creat test file")
	}
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %s", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.ChmodRequest{
		Path: "chmod.txt",
		Mode: 0777,
	}
	_, err = client.Chmod(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Chmod did not return OK. Actual: %s", err)
	}
}

func TestLink(t *testing.T) {
	// Create file to run test on
	_, err := commands.CreatCommand(tmpdir, "link.txt", os.FileMode(0777), false)
	if err != nil {
		t.Fatalf("Could not creat test file")
	}
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %s", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.LinkRequest{
		OldPath: "link.txt",
		NewPath: "linknew.txt",
	}
	_, err = client.Link(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Link did not return OK. Actual: %s", err)
	}
}

func TestPing(t *testing.T) {
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %s", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.PingRequest{
		Ip:   "0.0.0.0",
		Port: "0",
	}
	_, err = client.Ping(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Ping did not return OK. Actual: %s", err)
	}
}

func TestRename(t *testing.T) {
	// Create file to run test on
	_, err := commands.CreatCommand(tmpdir, "rename.txt", os.FileMode(0777), false)
	if err != nil {
		t.Fatalf("Could not creat test file")
	}
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %s", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.RenameRequest{
		OldPath: "rename.txt",
		NewPath: "renamenew.txt",
	}
	_, err = client.Rename(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Errorf("Rename did not return OK. Actual: %s", err)
	}
}

func TestTruncate(t *testing.T) {
	// Create file to run test on
	_, err := commands.CreatCommand(tmpdir, "truncate.txt", os.FileMode(0777), false)
	if err != nil {
		t.Fatalf("Could not creat test file")
	}
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Could not connect to server: %s", err)
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
	_, err := commands.CreatCommand(tmpdir, "unlink.txt", os.FileMode(0777), false)
	if err != nil {
		t.Fatalf("Could not creat test file")
	}
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
	_, err := commands.CreatCommand(tmpdir, "utimes.txt", os.FileMode(0777), false)
	if err != nil {
		t.Fatalf("Could not creat test file")
	}
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatal("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.UtimesRequest{
		Path:              "utimes.txt",
		AccessSeconds:     1,
		AccessNanoseconds: 1,
		ModifySeconds:     1,
		ModifyNanoseconds: 1,
	}
	_, err = client.Utimes(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Error("Utimes did not return OK. Actual:", err)
	}
}

func TestWrite(t *testing.T) {
	// Create file to run test on
	_, err := commands.CreatCommand(tmpdir, "write.txt", os.FileMode(0777), false)
	if err != nil {
		t.Fatalf("Could not creat test file")
	}
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

func TestMkdir(t *testing.T) {
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatal("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.MkdirRequest{
		Directory: "somedir",
		Mode:      0777,
	}
	_, err = client.Mkdir(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Error("Mkdir did not return OK. Actual:", err)
	}
}

func TestRmdir(t *testing.T) {
	// Create file to run test on
	_, err := commands.CreatCommand(tmpdir, "somedir.txt", os.FileMode(0777), false)
	if err != nil {
		t.Fatalf("Could not creat test file")
	}
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())
	if err != nil {
		t.Fatal("Could not connect to server:", err)
	}
	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	req := pb.RmdirRequest{
		Directory: "somedir",
	}
	_, err = client.Rmdir(context.Background(), &req)
	if grpc.Code(err) != codes.OK {
		t.Error("Rmdir did not return OK. Actual:", err)
	}
}

func createTestDir(t *testing.T, name string) {
	os.RemoveAll(path.Join(os.TempDir(), name))
	err := os.Mkdir(path.Join(os.TempDir(), name), 0777)
	if err != nil {
		t.Fatal("Error creating directory", err)
	}
}

func removeTestDir(name string) {
	time.Sleep(1 * time.Second)
	os.RemoveAll(path.Join(os.TempDir(), name))
}

func TestKillSignal(t *testing.T) {
	createTestDir(t, "testksMountpoint")
	defer removeTestDir("testksMountpoint")
	createTestDir(t, "testksDirectory")
	defer removeTestDir("testksDirectory")

	discovery := exec.Command("discovery-server", "--port=10102")
	err := discovery.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := discovery.Process.Kill()
		if err != nil {
			t.Error("Failed to kill discovery server,", err)
		}
	}()

	_, err = commands.InitCommand(path.Join(os.TempDir(), "testksDirectory"))
	if err != nil {
		t.Fatal(err)
	}

	pfsd := exec.Command("pfsd", path.Join(os.TempDir(), "testksDirectory"), path.Join(os.TempDir(), "testksMountpoint"), "localhost", "10102", "testPool")
	err = pfsd.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer pfsd.Process.Kill()

	time.Sleep(5 * time.Second)

	pidPath := path.Join(os.TempDir(), "testksDirectory", "meta", "pfsd.pid")
	if _, err := os.Stat(pidPath); err == nil {
		pidByte, err := ioutil.ReadFile(pidPath)
		if err != nil {
			t.Fatal("Can't read pid file", err)
		}
		pid, err := strconv.Atoi(string(pidByte))
		if err != nil {
			t.Fatal("Incorrect pid information", err)
		}
		err = syscall.Kill(pid, syscall.SIGTERM)
		if err != nil {
			t.Fatal("Can not kill PFSD,", err)
		}

		done := make(chan bool, 1)
		go func() {
			pfsd.Wait()
			done <- true
		}()

		select {
		case <-time.After(10 * time.Second):
			t.Fatal("pfsd did not finish within 10 seconds")
		case <-done:
			break
		}
	} else {
		t.Fatal("Could not read pid file:", err)
	}
}

func BenchmarkPing(b *testing.B) {
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	for n := 0; n < b.N; n++ {
		req := pb.PingRequest{
			Ip:   "0.0.0.0",
			Port: "0",
		}
		client.Ping(context.Background(), &req)
	}
}

func BenchmarkCreat(b *testing.B) {
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
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
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
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
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
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
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
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
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
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
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
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
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
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
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
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
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	defer conn.Close()
	conn, _ := grpc.Dial("localhost:10101", grpc.WithInsecure())
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		commands.MkdirCommand(tmpdir, "rmDirBench"+str, os.FileMode(0777), false)
		req := pb.RmdirRequest{
			Directory: "rmDirBench" + str,
		}
		client.Rmdir(context.Background(), &req)
	}
}
