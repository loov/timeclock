package pgdb

import (
	"github.com/loov/timeclock/project"
)

type Projects struct {
	*Database
}

func (db *Database) Projects() project.Database { return &Projects{db} }

func (db *Projects) CreateProject(project project.Project) (project.ID, error) {
	return 0, todo
}

func (db *Projects) Infos() ([]project.Info, error) {
	rows, err := db.Query(`
		SELECT ID, CustomerID, Slug, Name, Completed
		FROM Projects
		ORDER BY Name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []project.Info
	for rows.Next() {
		var p project.Info
		err := rows.Scan(
			&p.ID, &p.CustomerID, &p.Slug, &p.Name, &p.Completed,
		)
		if err != nil {
			return projects, err
		}
		projects = append(projects, p)
	}

	return projects, rows.Err()
}
