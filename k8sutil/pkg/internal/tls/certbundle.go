package tls

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

type CertificateBundle struct {
	AuthorityPEM         []byte
	BundlePEM            []byte
	ServerCertificatePEM []byte
	ServerKeyPEM         []byte
}

type CertificateBundleOptions struct {
	Organization        string
	CACommonName        string
	CertificateDNSNames []string
	Validity            time.Duration
}

type CertificateBundleOption interface {
	ApplyToCertificateBundleOptions(target *CertificateBundleOptions)
}

func (o *CertificateBundleOptions) ApplyOptions(opts []CertificateBundleOption) {
	for _, opt := range opts {
		opt.ApplyToCertificateBundleOptions(o)
	}
}

type CertificateBundleWithOrganization string

var _ CertificateBundleOption = CertificateBundleWithOrganization("")

func (wo CertificateBundleWithOrganization) ApplyToCertificateBundleOptions(target *CertificateBundleOptions) {
	target.Organization = string(wo)
}

type CertificateBundleWithCACommonName string

var _ CertificateBundleOption = CertificateBundleWithCACommonName("")

func (wccn CertificateBundleWithCACommonName) ApplyToCertificateBundleOptions(target *CertificateBundleOptions) {
	target.CACommonName = string(wccn)
}

type CertificateBundleWithCertificateDNSNames []string

var _ CertificateBundleOption = CertificateBundleWithCertificateDNSNames(nil)

func (wcdn CertificateBundleWithCertificateDNSNames) ApplyToCertificateBundleOptions(target *CertificateBundleOptions) {
	target.CertificateDNSNames = append(target.CertificateDNSNames, wcdn...)
}

type CertificateBundleWithValidity time.Duration

var _ CertificateBundleOption = CertificateBundleWithValidity(0)

func (wv CertificateBundleWithValidity) ApplyToCertificateBundleOptions(target *CertificateBundleOptions) {
	target.Validity = time.Duration(wv)
}

func GenerateCertificateBundle(opts ...CertificateBundleOption) (*CertificateBundle, error) {
	o := &CertificateBundleOptions{
		Organization: "Internet Widgits Pty Ltd",
		CACommonName: "Managed CA",
		Validity:     24 * time.Hour,
	}
	o.ApplyOptions(opts)

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	ca, err := cert.NewSelfSignedCACert(cert.Config{
		Organization: []string{o.Organization},
		CommonName:   o.CACommonName,
	}, caKey)
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
		Subject:      pkix.Name{Organization: []string{"Puppet"}},
		DNSNames:     o.CertificateDNSNames,
		SerialNumber: big.NewInt(1),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		NotBefore:    now,
		NotAfter:     now.Add(o.Validity),
	}
	if len(o.CertificateDNSNames) > 0 {
		serverCertTemplate.Subject.CommonName = o.CertificateDNSNames[0]
	}

	serverCertRaw, err := x509.CreateCertificate(rand.Reader, serverCertTemplate, ca, serverKey.Public(), caKey)
	if err != nil {
		return nil, err
	}

	cb := &CertificateBundle{
		AuthorityPEM:         pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.Raw}),
		ServerCertificatePEM: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCertRaw}),
		ServerKeyPEM:         pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: serverKeyRaw}),
	}
	cb.BundlePEM = append(append([]byte{}, cb.ServerCertificatePEM...), cb.AuthorityPEM...)
	return cb, nil
}
