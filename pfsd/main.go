package main

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Print("Usage:\n\tpfsd <port> <pfs_root_path>\n")
		os.Exit(1)
	}
	port, err := strconv.Atoi(os.Args[1])
	if err != nil || port < 1 || port > 65535 {
		log.Fatalln("FATAL: port must be a number between 1 and 65535, inclusive.")
	}
	pnetserver.PFSDir = os.Args[2]
	if _, err := os.Stat(pnetserver.PFSDir); os.IsNotExist(err) {
		log.Fatalln("FATAL: path", pnetserver.PFSDir, "does not exist.")
	}
	if _, err := os.Stat(path.Join(pnetserver.PFSDir, "meta")); os.IsNotExist(err) {
		log.Fatalln("FATAL: path", pnetserver.PFSDir, "is not valid PFS root.")
	}

	lis, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil {
		log.Fatalf("FATAL: Failed to listen on port %d: %v.\n", port, err)
	}
	srv := grpc.NewServer()
	pb.RegisterParanoidNetworkServer(srv, &pnetserver.ParanoidServer{})
	srv.Serve(lis)
}
