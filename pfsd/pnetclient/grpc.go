package pnetclient

import (
	"crypto/tls"
	"github.com/cpssd/paranoid/pfsd/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

var SkipVerify bool

func Dial(ipAddress globals.Node) *grpc.ClientConn {
	var opts []grpc.DialOption
	creds := credentials.NewTLS(&tls.Config{
		ServerName:         "",
		InsecureSkipVerify: SkipVerify,
	})
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(ipAddress.String(), opts...)
	if err != nil {
		log.Println("fail to dial: ", err)
	}
	return conn
}
