package main

import (
	"flag"
	"github.com/cpssd/paranoid/discovery-server/dnetserver"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strconv"
	"time"
)

var (
	port          = flag.Int("port", 10101, "port to listen on")
	logDir        = flag.String("log-directory", "/var/log", "directory in which to create ParanoidDiscovery.log")
	renewInterval = flag.Int("renew-interval", 5*60*1000, "time after which membership expires, in ms")
	certFile      = flag.String("cert", "", "TLS certificate file - if empty connection will be unencrypted")
	keyFile       = flag.String("key", "", "TLS key file - if empty connection will be unencrypted")
)

func createRPCServer() *grpc.Server {
	var opts []grpc.ServerOption
	if *certFile != "" && *keyFile != "" {
		dnetserver.Log.Info("Starting discovery server with TLS.")
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			dnetserver.Log.Fatal("Failed to generate TLS credentials:", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		dnetserver.Log.Info("Starting discovery server without TLS.")
	}
	return grpc.NewServer(opts...)
}

func main() {
	flag.Parse()
	var err error
	dnetserver.Log, err = logger.New("main", "discovery-server", *logDir)
	if err != nil {
		// TODO: Replace this with a call to a static logger in voy's package,
		// once it is implemented.
		log.Fatalln("Failed to create new logger:", err)
	}

	renewDuration, err := time.ParseDuration(strconv.Itoa(*renewInterval) + "ms")
	if err != nil {
		dnetserver.Log.Error("Failed parsing renew interval", err)
	}

	dnetserver.RenewInterval = renewDuration

	if *port < 1 || *port > 65535 {
		dnetserver.Log.Fatal("Port must be a number between 1 and 65535, inclusive.")
	}

	dnetserver.Log.SetOutput(logger.LOGFILE | logger.STDERR)
	if err != nil {
		dnetserver.Log.Fatal("Failed to set logger output:", err)
	}

	dnetserver.Log.Info("Starting Paranoid Discovery Server")

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		dnetserver.Log.Fatalf("Failed to listen on port %d: %v.", *port, err)
	}
	dnetserver.Log.Info("Listening on port", *port)

	srv := createRPCServer()
	pb.RegisterDiscoveryNetworkServer(srv, &dnetserver.DiscoveryServer{})
	go srv.Serve(lis)
	defer srv.Stop()
	dnetserver.Log.Info("gRPC server created")

	markInactiveNodes()
}

// Mark the nodes as inactive if their time expires
func markInactiveNodes() {
	for {
		now := time.Now()
		for i, node := range dnetserver.Nodes {
			dnetserver.Nodes[i].Active = node.ExpiryTime.Sub(now) > 0
		}
		time.Sleep(time.Second * 10)
	}
}
