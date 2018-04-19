package project

type Database interface {
	Infos() ([]Info, error)
}

type ID uint64

type Project struct {
	ID         ID
	CustomerID CustomerID

	Slug string
	Name string

	Activities  []string
	Description string
}

type Info struct {
	ID         ID
	CustomerID CustomerID

	Slug string
	Name string
}

type CustomerID uint64

type Customer struct {
	ID   CustomerID
	Slug string
	Name string
}
