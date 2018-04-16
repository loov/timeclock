package pgdb

import (
	"database/sql"
	"fmt"

	//TODO: switch to pgx
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func New(params string) (*Database, error) {
	sdb, err := sql.Open("postgres", params)
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

					Root  BOOL NOT NULL DEFAULT false,

					Created_At  TIMESTAMP NOT NULL DEFAULT current_timestamp,
					Modified_At TIMESTAMP NOT NULL DEFAULT current_timestamp
			)`,
			1: `CREATE TABLE Roles (
					UserID     INT8 PRIMARY KEY,
					
					Admin      BOOL NOT NULL DEFAULT false,
					Accountant BOOL NOT NULL DEFAULT false,
					Worker     BOOL NOT NULL DEFAULT false,

					Created_At  TIMESTAMP NOT NULL DEFAULT current_timestamp,
					Modified_At TIMESTAMP NOT NULL DEFAULT current_timestamp,

					FOREIGN KEY (UserID) REFERENCES Users (ID)
			)`,
		},
	}, {
		Name:    "Authentication",
		Version: 2,
		Steps: []string{
			0: `CREATE TABLE Credentials (
					UserID   INT8  NOT NULL,

					Provider TEXT  NOT NULL,
					Key      BYTEA NOT NULL,

					Created_At  TIMESTAMP NOT NULL DEFAULT current_timestamp,
					Modified_At TIMESTAMP NOT NULL DEFAULT current_timestamp,

					FOREIGN KEY (UserID) REFERENCES Users (ID)
			)`,
		},
	}, {
		Name:    "Customers",
		Version: 3,
		Steps: []string{
			0: `CREATE TABLE Customers (
					ID   BIGSERIAL PRIMARY KEY,
					Slug TEXT NOT NULL UNIQUE,
					Name TEXT NOT NULL,

					Created_At  TIMESTAMP NOT NULL DEFAULT current_timestamp,
					Modified_At TIMESTAMP NOT NULL DEFAULT current_timestamp
			)`,
		},
	}, {
		Name:    "Projects",
		Version: 4,
		Steps: []string{
			0: `CREATE TABLE Projects (
					ID          BIGSERIAL PRIMARY KEY,
					CustomerID  INT8,

					Slug TEXT NOT NULL,
					Name TEXT NOT NULL,
					Activities TEXT[] NOT NULL DEFAULT '{}',
					
					Description TEXT NOT NULL,

					Created_At  TIMESTAMP NOT NULL DEFAULT current_timestamp,
					Modified_At TIMESTAMP NOT NULL DEFAULT current_timestamp,

					FOREIGN KEY (CustomerID) REFERENCES Customers (ID)
			)`,
			1: `CREATE TABLE Activities (
					ID BIGSERIAL PRIMARY KEY,
					
					WorkerID  INT8 NOT NULL,
					ProjectID INT8 NOT NULL,

					Time     TIMESTAMP NOT NULL,
					Name     TEXT      NOT NULL,
					Amount   NUMERIC   NOT NULL,

					Locked   BOOL      NOT NULL DEFAULT false,
					
					Created_At  TIMESTAMP NOT NULL DEFAULT current_timestamp,
					Modified_At TIMESTAMP NOT NULL DEFAULT current_timestamp,

					FOREIGN KEY (WorkerID)  REFERENCES Users (ID),
					FOREIGN KEY (ProjectID) REFERENCES Projects (ID)
			)`,
		},
	},
}
