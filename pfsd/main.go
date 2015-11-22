package main

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
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
	if len(os.Args) < 5 {
		fmt.Print("Usage:\n\tpfsd <port> <paranoid_directory> <Discovery Server> <Discovery Port>\n")
		os.Exit(1)
	}
	port, err := strconv.Atoi(os.Args[1])
	if err != nil || port < 1 || port > 65535 {
		log.Fatalln("FATAL: port must be a number between 1 and 65535, inclusive.")
	}
	discoveryPort, err := strconv.Atoi(os.Args[4])
	if err != nil || discoveryPort < 1 || discoveryPort > 65535 {
		log.Fatalln("FATAL: Discovery port must be a number between 1 and 65535, inclusive.")
	}
	pnetclient.SetDiscovery(os.Args[3], os.Args[4], os.Args[1])
	pnetserver.ParanoidDir = os.Args[2]
	globals.Server, err = pnetclient.GetIP()
	if err != nil {
		log.Fatalln("FATAL: Cant get internal IP.")
	}
	globals.Port = port
	if _, err := os.Stat(pnetserver.ParanoidDir); os.IsNotExist(err) {
		log.Fatalln("FATAL: path", pnetserver.ParanoidDir, "does not exist.")
	}
	if _, err := os.Stat(path.Join(pnetserver.ParanoidDir, "meta")); os.IsNotExist(err) {
		log.Fatalln("FATAL: path", pnetserver.ParanoidDir, "is not valid PFS root.")
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("FATAL: Failed to listen on port %d: %v.\n", port, err)
	}

	pnetclient.JoinDiscovery("_")
	srv := grpc.NewServer()
	pb.RegisterParanoidNetworkServer(srv, &pnetserver.ParanoidServer{})
	srv.Serve(lis)
}
