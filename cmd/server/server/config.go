package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"
)

const DefaultHost = "localhost"
const DefaultGitAddress = DefaultHost + ":9999"
const DefaultWebAddress = DefaultHost + ":9998"
const DefaultCertPath = "data/server/server.crt"
const DefaultPrivateKey = "data/server/server.key"
const DefaultRepositoryPath = "data/server/repositories"
const DefaultDatabasePath = "data/server/db"

type WebConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type GitConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Config struct {
	GitConfig
	WebConfig

	Host       string
	CertPath   string
	PrivateKey string

	// DatabasePath points to a location where the database data is located. It can be a path
	// on the hard-drive
	DatabasePath string

	// RepositoryPath points to where repositories are located
	RepositoryPath string

	SuperUsername  string
	SuperPassword  string
	SuperPublicKey string
}

func (c *Config) ToTLSConfig() *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool
	caCert, err = ioutil.ReadFile(c.CertPath)
	if err != nil {
		log.Fatal("Error opening cert file", c.CertPath, ", error ", err)
	}
	caCertPool = x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		ServerName: c.Host,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12,
	}
}

func LoadConfig() Config {
	return Config{
		GitConfig: GitConfig{
			Address:      DefaultGitAddress,
			ReadTimeout:  5000,
			WriteTimeout: 5000,
			IdleTimeout:  5000,
		},
		WebConfig: WebConfig{
			Address:      DefaultWebAddress,
			ReadTimeout:  5000,
			WriteTimeout: 5000,
			IdleTimeout:  5000,
		},
		Host:           DefaultHost,
		CertPath:       DefaultCertPath,
		PrivateKey:     DefaultPrivateKey,
		RepositoryPath: DefaultRepositoryPath,
		DatabasePath:   DefaultDatabasePath,
	}
}
