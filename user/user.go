package user

import "strings"

type Database interface {
	Workers() ([]User, error)
}

type ID uint64

type User struct {
	ID ID

	Alias string
	Name  string
	Email string

	Root bool

	Roles
}

type Roles struct {
	Admin      bool
	Accountant bool
	Worker     bool
}

func (r Roles) String() string {
	var xs []string
	if r.Admin {
		xs = append(xs, "admin")
	}
	if r.Accountant {
		xs = append(xs, "accountant")
	}
	if r.Worker {
		xs = append(xs, "worker")
	}
	return strings.Join(xs, ", ")
}
