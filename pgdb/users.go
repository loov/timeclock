package pgdb

import "github.com/loov/timeclock/user"

type Users struct {
	*Database
}

func (users *Users) Register(user user.User, provider string, key []byte) (user.ID, error) {
	return 0, todo
}

func (users *Users) FindCredentials(name string) error {
	return todo
}

func (users *Users) List() ([]user.User, error) {
	return nil, todo
}
