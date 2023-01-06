package config

import (
	"crypto/tls"
)

// ServiceConfig represents the configuration of application's service
// and its security
type ServiceConfig struct {
	Port        int    `env:"PORT" json:"port"`
	SecuredPort int    `env:"SECURED_PORT" json:"securedPort"`
	TLSEnabled  bool   `env:"TLS_ENABLED" envDefault:"false" json:"tlsEnabled"`
	TLSCertPath string `env:"TLS_CERT_PATH" json:"tlsCertPath"`
	TLSKeyPath  string `env:"TLS_KEY_PATH" json:"tlsKeyPath"`

	tlsCert tls.Certificate
}

// LoadCerts attempts to load TLS certificates defined in configuration
func (c *ServiceConfig) LoadCerts() error {
	if !c.TLSEnabled {
		return nil
	}

	tlsCert, err := tls.LoadX509KeyPair(c.TLSCertPath, c.TLSKeyPath)
	if err != nil {
		return err
	}

	c.tlsCert = tlsCert

	return nil
}

// GetTLSConfig returns the *tls.Config based on given configuration
func (c *ServiceConfig) GetTLSConfig() *tls.Config {
	if !c.TLSEnabled {
		return nil
	}

	return &tls.Config{
		Certificates: []tls.Certificate{c.tlsCert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}
}
