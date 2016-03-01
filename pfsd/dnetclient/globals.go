package dnetclient

import (
	"crypto/tls"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

const peerPingTimeOut time.Duration = time.Minute * 3

var (
	discoveryCommonName string

	Log *logger.ParanoidLogger
)

func dialDiscovery() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(2*time.Second))
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
