package server

import (
	"golang.org/x/crypto/ssh"
	"os"
	"path"
	"strings"
	"unicode"
)

func testSpace(r rune) bool {
	return unicode.IsSpace(r) || r == '\u0005'
}

// SplitOnSpace the supplied string into two parts. The first part contains the command to
// be executed. The second part contains the arguments sent to the command. It will search for the first space and use
// that as the delimiter.
func SplitOnSpace(s string) []string {
	l := len(s)
	i := 0
	for i = 0; i < l; i++ {
		r := rune(s[i])
		if testSpace(r) {
			break
		}
	}
	if i == l {
		return []string{s}
	}
	result := []string{
		s[0:i],
		s[i+1:],
	}
	return result
}

// ReadPayload from the supplied request
func ReadPayload(req *ssh.Request) string {
	result := string(req.Payload)
	result = strings.Replace(result, "\x00", "", -1)
	l := len(result)
	for i := 0; i < l; i++ {
		v := rune(result[i])
		if unicode.IsLetter(v) {
			result = result[i:]
			break
		}
	}
	return result
}

// RepositoryExists checks if the supplied filename is a git repository
func RepositoryExists(filename string) bool {
	_, err := os.Stat(path.Join(filename, "HEAD"))
	return err == nil
}
