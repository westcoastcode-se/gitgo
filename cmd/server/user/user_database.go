package user

import "gitgo/api"

type Database interface {
	// GetUserUsingPublicKey can be called to search for potential users using a specific public key
	GetUserUsingPublicKey(publicKey string) *User
}

type DatabaseImpl struct {
	users []*User
}

func (d *DatabaseImpl) GetUserUsingPublicKey(fingerprint string) *User {
	for _, user := range d.users {
		for _, key := range user.PublicKeys {
			if key.Fingerprint == fingerprint {
				return user
			}
		}
	}
	return nil
}

func New() Database {
	return &DatabaseImpl{
		users: []*User{
			{
				Name:       "superuser",
				Password:   "superuser",
				PublicKeys: []api.PublicKey{},
			},
			{
				Name:     "per",
				Password: "password",
				PublicKeys: []api.PublicKey{
					{
						Name:        "MacOSX",
						Fingerprint: "SHA256:MWhJKWXhJL635pzRPGjovQbrLZcTpwB7g2NWdifQKIQ",
						PublicKey:   "",
					},
				},
			},
		},
	}
}
