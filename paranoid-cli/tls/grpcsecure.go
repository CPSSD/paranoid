package tls

import (
	"crypto/tls"
	"google.golang.org/grpc/credentials"
	"os"
	"path"
)

func CertExists(pfsDir string) bool {
	certPath := path.Join(pfsDir, "meta", "cert.pem")
	keyPath := path.Join(pfsDir, "meta", "key.pem")
	return exists(certPath) && exists(keypath)
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GenCreds() credentials.TransportAuthenticator {
	creds := credentials.NewTLS(&tls.Config{
		ServerName:         discoveryCommonName,
		InsecureSkipVerify: globals.TLSSkipVerify,
	})
	return creds
}
