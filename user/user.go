package user

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
