package user

import "gitgo/api"

type User struct {
	Name       string
	Password   string
	PublicKeys []api.PublicKey
}

func (u *User) ToApi() *api.User {
	return &api.User{
		Name:       u.Name,
		Password:   u.Password,
		PublicKeys: u.PublicKeys,
	}
}
