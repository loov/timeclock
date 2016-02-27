package user

type ID int64

type User struct {
	ID    ID
	Name  string
	Email string

	Accountant bool
	Worker     bool
	Reviewer   bool
	Supervisor bool
	Inactive   bool
}
