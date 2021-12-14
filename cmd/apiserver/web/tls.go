package web

import (
	"crypto/tls"
)

// TryExtractCommonName will try to extract the client certificates common name. This is assumed to be the username
func TryExtractCommonName(connectionState *tls.ConnectionState) string {
	if connectionState != nil && len(connectionState.VerifiedChains) > 0 && len(connectionState.VerifiedChains[0]) > 0 {
		var commonName = connectionState.VerifiedChains[0][0].Subject.CommonName
		return commonName
	}
	return ""
}
