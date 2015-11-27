package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func Dial(ipAddress globals.Node) *grpc.ClientConn {
	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(nil, "")
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, opts...)
	if err != nil {
		log.Println("fail to dial: ", err)
	}
	return conn
}
