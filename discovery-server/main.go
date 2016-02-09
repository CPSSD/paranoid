package main

import (
	"flag"
	"github.com/cpssd/paranoid/discovery-server/dnetserver"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	dnetserver.Log = logger.New("main", "discovery-server", *logDir)
	err := dnetserver.Log.SetOutput(logger.LOGFILE | logger.STDERR)
	if err != nil {
		dnetserver.Log.Error("Failed to set logger output:", err)
	}

	renewDuration, err := time.ParseDuration(strconv.Itoa(*renewInterval) + "ms")
	if err != nil {
		dnetserver.Log.Error("Failed parsing renew interval", err)
	}

	dnetserver.RenewInterval = renewDuration

	if *port < 1 || *port > 65535 {
		dnetserver.Log.Fatal("Port must be a number between 1 and 65535, inclusive.")
	}

	dnetserver.Log.Info("Starting Paranoid Discovery Server")

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		dnetserver.Log.Fatalf("Failed to listen on port %d: %v.", *port, err)
	}
	dnetserver.Log.Info("Listening on port", *port)

	srv := createRPCServer()
	pb.RegisterDiscoveryNetworkServer(srv, &dnetserver.DiscoveryServer{})

	dnetserver.Log.Info("gRPC server created")
	srv.Serve(lis)
}
