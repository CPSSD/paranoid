package main

import (
	"flag"
	"fmt"
	"github.com/cpssd/paranoid/network/discovery-server/dnetserver"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	var port int
	var logDir string
	var renewInterval int

	flag.IntVar(&port, "-port", 10101, "Server Port")
	flag.StringVar(&logDir, "-log-directory", "/var/log", "Log Directory")
	flag.IntVar(&renewInterval, "-renew-interval", 5*60*1000, "Time after membership expires, in ms")

	flag.Parse()

	renewDuration, _ := time.ParseDuration(strconv.Itoa(renewInterval) + "ms")

	dnetserver.RenewInterval = renewDuration

	if port < 1 || port > 65535 {
		fmt.Println("FATAL: port must be a number between 1 and 65535, inclusive.")
		os.Exit(1)
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		fmt.Println("FATAL: Log path", logDir, "does not exist.")
		os.Exit(1)
	}

	logFilePath := path.Join(logDir, "ParanoidDiscovery.log")
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("FATAL: Cannot write to file", logFilePath)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Println("[I] Starting Paranoid Discovery Server...")

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("[F] Failed to listen on port %d: %v.\n", port, err)
	}
	log.Println("[I] Listening on port", port)

	srv := grpc.NewServer()
	pb.RegisterDiscoveryNetworkServer(srv, &dnetserver.DiscoveryServer{})
	srv.Serve(lis)
	log.Println("[I] gRPC server created")

	go removeLoop()
}

// Mark the nodes as inactive if their time expires
func removeLoop() {
	for {
		now := time.Now()
		for i, node := range dnetserver.Nodes {
			dnetserver.Nodes[i].Active = node.ExpiryTime.Sub(now) < 0
		}
	}
}
