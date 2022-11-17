package config_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Icikowski/GoosyMock/config"
	"github.com/stretchr/testify/require"
)

func noopCertGen(t *testing.T) (tlsCert string, tlsKey string) {
	t.Helper()
	return
}

func allOkCertGen(t *testing.T) (tlsCert string, tlsKey string) {
	t.Helper()

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(123456789),
		Subject: pkix.Name{
			Country:            []string{"PL"},
			Province:           []string{"mazowieckie"},
			Locality:           []string{"Warszawa"},
			PostalCode:         []string{"00-000"},
			StreetAddress:      []string{"Testowa 13"},
			Organization:       []string{"Firma"},
			OrganizationalUnit: []string{"Certificate Authority"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Country:            []string{"PL"},
			Province:           []string{"mazowieckie"},
			Locality:           []string{"Warszawa"},
			PostalCode:         []string{"00-000"},
			StreetAddress:      []string{"Testowa 13"},
			Organization:       []string{"Firma"},
			OrganizationalUnit: []string{"Server"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(5, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	certFile, err := os.CreateTemp(os.TempDir(), "tls-*.crt")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Cleanup(func() {
		os.Remove(certFile.Name())
	})

	if err := pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		t.Fatal(err.Error())
		return
	}
	tlsCert = certFile.Name()

	keyFile, err := os.CreateTemp(os.TempDir(), "tls-*.key")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Cleanup(func() {
		os.Remove(keyFile.Name())
	})

	if err := pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	}); err != nil {
		t.Fatal(err.Error())
		return
	}
	tlsKey = keyFile.Name()

	return
}

func TestLoadCerts(t *testing.T) {
	tests := map[string]struct {
		sslEnabled    bool
		genFunc       func(*testing.T) (string, string)
		errorExpected bool
	}{
		"SSL disabled": {
			sslEnabled:    false,
			genFunc:       noopCertGen,
			errorExpected: false,
		},
		"TLS missing": {
			sslEnabled:    true,
			genFunc:       noopCertGen,
			errorExpected: true,
		},
		"TLS valid": {
			sslEnabled:    true,
			genFunc:       allOkCertGen,
			errorExpected: false,
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			conf := &config.ServiceConfig{
				SSLEnabled: tc.sslEnabled,
			}
			conf.TLSCertPath, conf.TLSKeyPath = tc.genFunc(t)

			err := conf.LoadCerts()
			if tc.errorExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetTLSConfig(t *testing.T) {
	tests := map[string]struct {
		sslEnabled bool
	}{
		"SSL disabled": {sslEnabled: false},
		"SSL enabled":  {sslEnabled: true},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			conf := &config.ServiceConfig{
				SSLEnabled: tc.sslEnabled,
			}
			actual := conf.GetTLSConfig()

			if !tc.sslEnabled {
				require.Nil(t, actual)
				return
			}
			require.NotNil(t, actual)
			require.EqualValues(t, tls.VersionTLS12, actual.MinVersion)
			require.ElementsMatch(t, []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			}, actual.CipherSuites)
		})
	}
}
