package project

type ID uint64

type Project struct {
	ID   ID
	Name string

	ShortDescription  string
	DefaultActivities []string
}
