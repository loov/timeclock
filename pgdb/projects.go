package pgdb

import "github.com/loov/timeclock/project"

type Projects struct {
	*Database
}

func (db *Projects) CreateProject(project project.Project) (project.ID, error) {
	return 0, todo
}
