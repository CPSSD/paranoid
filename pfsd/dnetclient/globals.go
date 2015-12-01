package dnetclient

import (
	"crypto/tls"
	"github.com/cpssd/paranoid/pfsd/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

const peerPingInterval time.Duration = time.Minute

var (
	// ThisNode has to be set before calling Join
	ThisNode            globals.Node
	discoveryCommonName string
)

func dialDiscovery() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if globals.TLSEnabled {
		creds := credentials.NewTLS(&tls.Config{
			ServerName:         discoveryCommonName,
			InsecureSkipVerify: globals.TLSSkipVerify,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.Dial(globals.DiscoveryAddr, opts...)
}
