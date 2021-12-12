package api

type User struct {
	// Name represents a unique name for a user
	Name string

	// Password is the user's password
	Password string

	// PublicKeys is all public keys for this user
	PublicKeys []PublicKey
}
