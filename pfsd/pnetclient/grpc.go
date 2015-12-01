package pnetclient

import (
	"crypto/tls"
	"github.com/cpssd/paranoid/pfsd/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func Dial(node globals.Node) *grpc.ClientConn {
	var opts []grpc.DialOption
	if globals.TLSEnabled {
		creds := credentials.NewTLS(&tls.Config{
			ServerName:         node.CommonName,
			InsecureSkipVerify: globals.TLSSkipVerify,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(node.String(), opts...)
	if err != nil {
		log.Println("fail to dial: ", err)
	}
	return conn
}
