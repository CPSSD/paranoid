package tls

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// GenCertificate will generate a TLS cert and key, prompting the user
// to enter relevant information. The resulting data will be saved to PEM
// files inside outputDir called "cert.pem" and "key.pem".
func GenCertificate(pfsDir string) error {
	scanner := bufio.NewScanner(os.Stdin)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("could not generate private key: %v", err)
	}

	startDate := time.Now()
	fmt.Print("Enter the length of time for which the cert will be valid, in days: ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not read user input: %v", err)
	}
	validForString := scanner.Text()
	validForDays, err := strconv.Atoi(validForString)
	if err != nil {
		return fmt.Errorf("could not parse user input as integer: %v", err)
	}
	validForString = strconv.Itoa(validForDays*24) + "h"
	validFor, err := time.ParseDuration(validForString)
	if err != nil {
		return fmt.Errorf("could not parse user input as duration: %v", err)
	}
	endDate := startDate.Add(validFor)

	// Sets max to 1 << 128
	serialNumberMax := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberMax)
	if err != nil {
		return fmt.Errorf("could not generate serial number: %v", err)
	}

	fmt.Print("Enter the name of your organisation: ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not read user input: %v", err)
	}
	organisation := scanner.Text()

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			// American spelling :(
			Organization: []string{organisation},
		},
		NotBefore:             startDate,
		NotAfter:              endDate,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}

	fmt.Print("Enter a comma-separated list of hostnames and/or IP addresses this cert will validate: ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not read user input: %v", err)
	}
	hostsString := scanner.Text()

	hosts := strings.Split(hostsString, ",")
	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	certPath := path.Join(pfsDir, "meta", "cert.pem")
	certFile, err := os.OpenFile(certPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer certFile.Close()
	if err != nil {
		return fmt.Errorf("failed to create cert file: %v", err)
	}
	pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	log.Println("INFO: Wrote certificate to", certPath)

	keyPath := path.Join(pfsDir, "meta", "key.pem")
	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer keyFile.Close()
	if err != nil {
		return fmt.Errorf("failed to create key file: %v", err)
	}
	pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	log.Println("INFO: Wrote key to", keyPath)

	return nil
}
