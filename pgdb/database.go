package pgdb

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
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

func (db *Database) DANGEROUS_DROP_ALL_TABLES() error {
	_, err := db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)
	return err
}

var migrations = Migrations{
	{
		Name:    "Users",
		Version: 1,
		Steps: []string{
			0: `CREATE TABLE Users (
				ID    BIGSERIAL PRIMARY KEY,
				Alias TEXT NOT NULL,
				Name  TEXT NOT NULL,
				Email TEXT NOT NULL,

				Root  BOOL NOT NULL DEFAULT false
			)`,
		},
	}, {
		Name:    "Authentication",
		Version: 2,
		Steps: []string{
			0: `CREATE TABLE Credentials (
					UserID   INT8,

					Provider TEXT  NOT NULL,
					Key      BYTEA NOT NULL,

					FOREIGN KEY (UserID) REFERENCES Users (ID)
			)`,
		},
	}, {
		Name:    "Companies",
		Version: 3,
		Steps: []string{
			0: `CREATE TABLE Companies (
					ID   BIGSERIAL PRIMARY KEY,
					Slug TEXT NOT NULL,
					Name TEXT NOT NULL
			)`,
			1: `CREATE TABLE Roles (
					UserID     INT8 PRIMARY KEY,
					CompanyID  INT8,
					
					Admin      BOOL NOT NULL DEFAULT false,
					Accountant BOOL NOT NULL DEFAULT false,
					Worker     BOOL NOT NULL DEFAULT false,
	
					FOREIGN KEY (UserID) REFERENCES Users (ID),
					FOREIGN KEY (CompanyID) REFERENCES Companies (ID)
			)`,
		},
	}, {
		Name:    "Projects",
		Version: 4,
		Steps: []string{
			0: `CREATE TABLE Projects (
					ID   BIGSERIAL PRIMARY KEY,
					Slug TEXT NOT NULL,
					Name TEXT NOT NULL,

					Activities TEXT[] NOT NULL DEFAULT '{}'
			)`,
			1: `CREATE TABLE Activities (
					ID BIGSERIAL PRIMARY KEY,
					
					WorkerID  INT8 NOT NULL,
					ProjectID INT8 NOT NULL,

					Time     TIMESTAMP NOT NULL,
					Name     TEXT      NOT NULL,
					Amount   NUMERIC   NOT NULL,

					Locked   BOOL      NOT NULL DEFAULT false,
					
					CreatedAt  TIMESTAMP NOT NULL DEFAULT current_timestamp,
					ModifiedAt TIMESTAMP NOT NULL DEFAULT current_timestamp,

					FOREIGN KEY (WorkerID)  REFERENCES Users (ID),
					FOREIGN KEY (ProjectID) REFERENCES Projects (ID)
			)`,
		},
	},
}
