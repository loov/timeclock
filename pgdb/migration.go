package pgdb

import (
	"database/sql"
	"fmt"
)

var versionTable = Migration{
	Name:    "Version Table",
	Version: 0,
	Steps: []string{
		0: `
		CREATE TABLE IF NOT EXISTS Versions (
			Version INT       NOT NULL UNIQUE,
			Updated TIMESTAMP NOT NULL DEFAULT current_timestamp
		)`,
	},
}

type Migrations []*Migration

func (migs Migrations) Run(db *Database) error {
	// check whether we have versions table
	// query := `SELECT COUNT(1) FROM information_schema.tables WHERE table_name = $1 AND table_schema = (SELECT current_schema()) LIMIT 1`
	// if err := p.db.QueryRow(query, "data_migrations").Scan(&count); err != nil {
	// 	return err
	// }

	var count int
	err := db.QueryRow(`SELECT COUNT(1) 
		FROM information_schema.tables
		WHERE table_name = lower($1) AND
		table_schema = (SELECT current_schema())`, "Versions").Scan(&count)
	if err == sql.ErrNoRows || count == 0 {
		fmt.Println(err, count)
		if err := versionTable.Run(db); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("migrations: unable to query versions table: %v", err)
	}

	// find last version
	version := 0
	err = db.QueryRow(`SELECT MAX(Version) FROM Versions`).Scan(&version)
	if err != nil {
		return fmt.Errorf("migrations: unable to query last version: %v", err)
	}

	// run all new migrations
	for _, mig := range migs {
		if mig.Version <= version {
			continue
		}

		if err := mig.Run(db); err != nil {
			return err
		}
	}

	return nil
}

type Migration struct {
	Name    string
	Version int
	Steps   []string
}

func (mig *Migration) Run(db *Database) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, step := range mig.Steps {
		_, err := tx.Exec(step)
		if err != nil {
			return fmt.Errorf("migration %q (v%v): failed at step %v: %v", mig.Name, mig.Version, i, err)
		}
	}

	_, err = tx.Exec(`INSERT INTO Versions (Version) VALUES ($1)`, mig.Version)
	if err != nil {
		return fmt.Errorf("migration %q (v%v): unable to update version: %v", mig.Name, mig.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("migration %q (v%v): unable to commit: %v", mig.Name, mig.Version, err)
	}

	return nil
}
