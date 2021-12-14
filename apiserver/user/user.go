package user

import "github.com/westcoastcode-se/gitgo/api"

type Users struct {
	Users []*User
}

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
