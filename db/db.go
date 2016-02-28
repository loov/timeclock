package db

type DB struct{}

func New(params string) (*DB, error) {
	return &DB{}, nil
}
