package server

import (
	"time"
)

const DefaultAddress = ":9998"
const DefaultCAPath = "data/ca.crt"
const DefaultCertPath = "data/apiserver.crt"
const DefaultPrivateKey = "data/apiserver.key"
const DefaultRepositoryPath = "data/repositories"
const DefaultDatabasePath = "data/db"

type Config struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	CAPath     string
	CertPath   string
	PrivateKey string

	// DatabasePath points to a location where the database data is located. It can be a path
	// on the hard-drive
	DatabasePath string

	// RepositoryPath points to where repositories are located
	RepositoryPath string
}

func LoadConfig() Config {
	return Config{
		Address:        DefaultAddress,
		ReadTimeout:    5000 * time.Millisecond,
		WriteTimeout:   5000 * time.Millisecond,
		IdleTimeout:    5 * time.Minute,
		CAPath:         DefaultCAPath,
		CertPath:       DefaultCertPath,
		PrivateKey:     DefaultPrivateKey,
		RepositoryPath: DefaultRepositoryPath,
		DatabasePath:   DefaultDatabasePath,
	}
}
