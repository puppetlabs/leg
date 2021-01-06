package tlsproxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"k8s.io/client-go/util/cert"
)

var tlsCAConfig = cert.Config{
	CommonName:   "Puppet Leg Kubernetes Integration Test CA",
	Organization: []string{"Puppet"},
}

type tlsCertificateBundle struct {
	AuthorityPEM         []byte
	BundlePEM            []byte
	ServerCertificatePEM []byte
	ServerKeyPEM         []byte
}

func generateTLSCertificateBundle(commonName string) (*tlsCertificateBundle, error) {
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	ca, err := cert.NewSelfSignedCACert(tlsCAConfig, caKey)
	if err != nil {
		return nil, err
	}

	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	serverKeyRaw, err := x509.MarshalPKCS8PrivateKey(serverKey)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	serverCertTemplate := &x509.Certificate{
		Subject:      pkix.Name{CommonName: commonName, Organization: []string{"Puppet"}},
		DNSNames:     []string{commonName},
		SerialNumber: big.NewInt(1),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		NotBefore:    now,
		NotAfter:     now.Add(24 * time.Hour),
	}

	serverCertRaw, err := x509.CreateCertificate(rand.Reader, serverCertTemplate, ca, serverKey.Public(), caKey)
	if err != nil {
		return nil, err
	}

	cb := &tlsCertificateBundle{
		AuthorityPEM:         pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.Raw}),
		ServerCertificatePEM: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCertRaw}),
		ServerKeyPEM:         pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: serverKeyRaw}),
	}
	cb.BundlePEM = append(append([]byte{}, cb.ServerCertificatePEM...), cb.AuthorityPEM...)
	return cb, nil
}
