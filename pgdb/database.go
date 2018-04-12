package pgdb

import (
	"database/sql"
	"fmt"
)

type Database struct {
	*sql.DB
}

func New(db, params string) (*Database, error) {
	sdb, err := sql.Open(db, params)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err)
	}

	return &Database{DB: sdb}, nil
}

func (db *Database) Init() error {
	return migrations.Run(db)
}

var migrations = Migrations{}
