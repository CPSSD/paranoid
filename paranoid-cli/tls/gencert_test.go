// +build !integration !benchmark

package tls

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/cpssd/paranoid/logger"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestGenerateCert(t *testing.T) {
	Log = logger.New("tlstest", "tlstest", os.DevNull)
	testPath := path.Join(os.TempDir(), "testCertGen")
	os.RemoveAll(testPath)
	os.Mkdir(testPath, 0777)
	os.Mkdir(path.Join(testPath, "meta"), 0777)
	// Since GenCertificate takes input from os.Stdin, we need to create
	// a fake stdin.
	fakeStdin, err := os.OpenFile(path.Join(testPath, "fakeStdin"),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	defer fakeStdin.Close()
	if err != nil {
		t.Fatal("Failed to create fake stdin file:", err)
	}
	_, err = fakeStdin.WriteString("365\nParanoid Inc.\ntest.paranoid.com\n")
	if err != nil {
		t.Fatal("Failed to write to fake stdin file:", err)
	}
	// Force a flush of the newly-written data.
	fakeStdin.Close()
	fakeStdin, err = os.Open(path.Join(testPath, "fakeStdin"))
	if err != nil {
		t.Fatal("Failed to reopen fake stdin file:", err)
	}
	os.Stdin = fakeStdin
	if err != nil {
		t.Fatal("Failed to redirect stdout:", err)
	}

	err = GenCertificate(testPath)
	if err != nil {
		t.Fatal("GenCertificate returned error:", err)
	}

	// Now we read the generated cert to confirm everything went okay.
	certPath := path.Join(testPath, "meta", "cert.pem")
	encodedPem, err := ioutil.ReadFile(certPath)
	if err != nil {
		t.Fatalf("could not open %s: %v\n", certPath, err)
	}
	decodedPem, _ := pem.Decode(encodedPem)
	if decodedPem == nil {
		t.Fatal("could not decode PEM file")
	}
	cert, err := x509.ParseCertificate(decodedPem.Bytes)
	if err != nil {
		t.Fatal("could not parse certificate:", err)
	}
	if cert.Subject.Organization[0] != "Paranoid Inc." {
		t.Error("Organisation field incorrect. Expected: \"Paranoid Inc.\". Actual:",
			cert.Subject.Organization[0])
	}
	certDuration := cert.NotAfter.Sub(cert.NotBefore)
	expectedDuration, _ := time.ParseDuration("8760h")
	if certDuration != expectedDuration {
		t.Errorf("Certificate duration incorrect. Expected: %s. Actual: %s\n",
			expectedDuration, certDuration)
	}
	if cert.DNSNames[0] != "test.paranoid.com" {
		t.Error("DNS name incorrect. Expected: \"test.paranoid.com\". Actual:",
			cert.DNSNames[0])
	}
}
