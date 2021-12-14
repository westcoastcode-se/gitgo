package user

import (
	"errors"
	"fmt"
	"gitgo/api"
	"gitgo/apiserver/db"
	"gitgo/apiserver/event"
	"sync"
)

const DatabasePath = "/users.json"

type Database interface {
	// AddUser adds the supplied user
	AddUser(user *User) error

	// GetUserUsingPublicKey can be called to search for potential users using a specific public key
	GetUserUsingPublicKey(publicKey string) *User
}

type DatabaseImpl struct {
	// Database is a generic json database
	contentDatabase db.ContentDatabase

	users []*User
	mutex *sync.RWMutex
}

func (d *DatabaseImpl) AddUser(newUser *User) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, user := range d.users {
		if user.Name == newUser.Name {
			return errors.New("user already exists")
		}
	}
	d.users = append(d.users, newUser)
	err := d.contentDatabase.Write(DatabasePath, &Users{d.users},
		fmt.Sprintf("adding user %s", newUser.Name))
	if err != nil {
		return err
	}
	// TODO post event on event queue
	return nil
}

func (d *DatabaseImpl) OnEvent(event event.Event) error {
	if e, ok := event.(*db.EventDataChanged); ok {
		if e.Path == DatabasePath {
			return d.reload()
		}
	}
	return nil
}

func (d *DatabaseImpl) GetUserUsingPublicKey(fingerprint string) *User {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, user := range d.users {
		for _, key := range user.PublicKeys {
			if key.Fingerprint == fingerprint {
				return user
			}
		}
	}
	return nil
}

func (d *DatabaseImpl) reload() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var users Users
	err := d.contentDatabase.Read(DatabasePath, &users)
	if err != nil {
		return err
	}
	d.users = users.Users
	return nil
}

func New(database db.ContentDatabase) Database {
	return &DatabaseImpl{
		contentDatabase: database,
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
		mutex: &sync.RWMutex{},
	}
}
