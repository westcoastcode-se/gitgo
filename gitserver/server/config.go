package server

import "time"

const DefaultAddress = ":9999"
const DefaultRepositoriesPath = "data/repositories"

type Config struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// GitBinDir points to the directory where git binaries are located
	GitBinDir string

	// Root path where the actual repositories are found
	RepositoriesPath string

	// Path to where the key used by the SSH server when exposing itself to
	// a client
	SSHKeyPath string

	// APIServerAddress is the address to the api server
	APIServerAddress string

	// The part to where a PEM encoded certificate file is located
	ClientCertPath string

	// The part to where a PEM encoded private key file is located
	ClientKeyPath string

	// CA used when issuing the client certificate
	ClientCAPath string

	// Should the server ignore if the keys are considered insecure. This normally happens when
	// you are using self-signed certificates
	InsecureSkipVerify bool
}

func LoadConfig() *Config {
	cfg := &Config{
		Address:            DefaultAddress,
		ReadTimeout:        5000 * time.Millisecond,
		WriteTimeout:       5000 * time.Millisecond,
		IdleTimeout:        5000 * time.Millisecond,
		GitBinDir:          "C:\\Program Files\\Git\\mingw64\\bin",
		RepositoriesPath:   DefaultRepositoriesPath,
		SSHKeyPath:         "data/gitserver.key",
		APIServerAddress:   "https://localhost:9998",
		ClientCertPath:     "data/apiserver_client.crt",
		ClientKeyPath:      "data/apiserver_client.key",
		ClientCAPath:       "data/ca.crt",
		InsecureSkipVerify: true,
	}
	return cfg
}
