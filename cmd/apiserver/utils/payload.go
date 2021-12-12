package utils

import (
	"golang.org/x/crypto/ssh"
	"strings"
	"unicode"
)

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
