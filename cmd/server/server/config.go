package server

const DefaultAddress = ":9999"
const DefaultPrivateKey = "data/server/server.key"
const DefaultRepositoryPath = "data/server/repositories"

type Config struct {
	Address    string
	PrivateKey string

	// RepositoryPath points to where repositories are located
	RepositoryPath string

	SuperUsername  string
	SuperPassword  string
	SuperPublicKey string
}

func LoadConfig() Config {
	return Config{
		Address:        DefaultAddress,
		PrivateKey:     DefaultPrivateKey,
		RepositoryPath: DefaultRepositoryPath,
	}
}
