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

var migrations = Migrations{
	{
		Name:    "Users",
		Version: 1,
		Steps: []string{
			0: `CREATE TABLE users.Users (
				ID    INT8 AUTO_INCREMENT PRIMARY KEY,
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
			0: `CREATE TABLE users.Credentials (
					UserID     INT8,

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
					ID   INT8 AUTO_INCREMENT PRIMARY KEY,
					Slug TEXT NOT NULL,
					Name TEXT NOT NULL
			)`,
			1: `CREATE TABLE Roles (
					UserID    INT8 NOT NULL,
					
					Admin      BOOL NOT NULL DEFAULT false,
					Accountant BOOL NOT NULL DEFAULT false,
					Worker     BOOL NOT NULL DEFAULT false,
	
					FOREIGN KEY (UserID) REFERENCES Users (ID)
			)`,
		},
	}, {
		Name:    "Projects",
		Version: 4,
		Steps: []string{
			0: `CREATE TABLE Projects (
					ID   INT8 AUTO_INCREMENT PRIMARY KEY,
					Slug TEXT NOT NULL,
					Name TEXT NOT NULL,

					Activities TEXT[] NOT NULL DEFAULT '{}'
			)`,
			1: `CREATE TABLE Activities (
					ID INT8 AUTO_INCREMENT PRIMARY KEY,
					
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
