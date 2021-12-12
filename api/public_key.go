package api

type PublicKey struct {
	// Name is a unique name of this public key for a specific user
	Name string

	// Fingerprint is the fingerprint of the public key
	Fingerprint string

	// PublicKey is the public key
	PublicKey string
}
