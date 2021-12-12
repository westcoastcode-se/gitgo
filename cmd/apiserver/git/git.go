package git

import (
	"os"
	"path"
)

func RepositoryExists(p string) bool {
	_, err := os.Stat(path.Join(p, "HEAD"))
	return err == nil
}
