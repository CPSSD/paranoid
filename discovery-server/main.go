package main

import (
	"flag"
	"fmt"
	"github.com/cpssd/paranoid/discovery-server/dnetserver"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	port          = flag.Int("port", 10101, "port to listen on")
	logDir        = flag.String("log_directory", "/var/log", "directory in which to create ParanoidDiscovery.log")
	renewInterval = flag.Int("renew-interval", 5*60*1000, "time after which membership expires, in ms")
	certFile      = flag.String("cert", "", "TLS certificate file - if empty connection will be unencrypted")
	keyFile       = flag.String("key", "", "TLS key file - if empty connection will be unencrypted")
)

func createRPCServer() *grpc.Server {
	var opts []grpc.ServerOption
	if *certFile != "" && *keyFile != "" {
		log.Println("INFO: Starting discovery server with TLS.")
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalln("FATAL: Failed to generate TLS credentials:", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		log.Println("INFO: Starting discovery server without TLS.")
	}
	return grpc.NewServer(opts...)
}

func main() {
	flag.Parse()

	renewDuration, _ := time.ParseDuration(strconv.Itoa(*renewInterval) + "ms")

	dnetserver.RenewInterval = renewDuration

	if *port < 1 || *port > 65535 {
		fmt.Println("FATAL: port must be a number between 1 and 65535, inclusive.")
		os.Exit(1)
	}

	if _, err := os.Stat(*logDir); os.IsNotExist(err) {
		fmt.Println("FATAL: Log path", *logDir, "does not exist.")
		os.Exit(1)
	}

	logFilePath := path.Join(*logDir, "ParanoidDiscovery.log")
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("FATAL: Cannot write to file", logFilePath)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Println("[I] Starting Paranoid Discovery Server...")

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("[F] Failed to listen on port %d: %v.\n", *port, err)
	}
	log.Println("[I] Listening on port", *port)

	srv := createRPCServer()
	pb.RegisterDiscoveryNetworkServer(srv, &dnetserver.DiscoveryServer{})
	go srv.Serve(lis)
	defer srv.Stop()
	log.Println("[I] gRPC server created")

	markInactiveNodes()
}

// Mark the nodes as inactive if their time expires
func markInactiveNodes() {
	for {
		now := time.Now()
		for i, node := range dnetserver.Nodes {
			dnetserver.Nodes[i].Active = node.ExpiryTime.Sub(now) < 0
		}
		time.Sleep(time.Second * 10)
	}
}
