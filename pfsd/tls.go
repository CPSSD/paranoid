// tls.go contains utility functions relating to TLS support

package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func getCommonNameFromCert(certFile string) (string, error) {
	encodedPem, err := ioutil.ReadFile(certFile)
	if err != nil {
		return "", fmt.Errorf("could not open %s: %v", certFile, err)
	}
	decodedPem, _ := pem.Decode(encodedPem)
	if decodedPem == nil {
		return "", fmt.Errorf("could not decode PEM file")
	}
	cert, err := x509.ParseCertificate(decodedPem.Bytes)
	if err != nil {
		return "", fmt.Errorf("could not parse certificate: %v", err)
	}
	return cert.Subject.CommonName, nil
}
