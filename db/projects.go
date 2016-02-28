package db

import (
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

func (db *DB) Projects() project.Projects {
	return &Projects{
		list: []project.Project{
			upgradeofwater,
			coversfloatation,
		},
	}
}

type Projects struct {
	list []project.Project
}

func (projects *Projects) List() ([]project.Project, error) {
	return projects.list, nil
}

func (projects *Projects) ByID(id project.ID) (project.Project, error) {
	for _, p := range projects.list {
		if p.ID == id {
			return p, nil
		}
	}

	return project.Project{ID: id}, project.ErrNotExist
}

var day = 24 * 60 * 60 * time.Second

var upgradeofwater = project.Project{
	ID:          "15113",
	Caption:     "Upgrade of watering system",
	Customer:    "ACME",
	Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Magnam impedit cumque nam necessitatibus quod hic possimus rerum, eveniet repudiandae! Ex quis unde provident, explicabo commodi ullam quibusdam enim officiis quaerat.",
	Status:      project.Active,

	Engineers: []user.ID{},
	Estimate:  30 * day,

	Created:   time.Now().Add(-day * 5),
	Modified:  time.Now().Add(-day * 2),
	Completed: time.Time{},
}

var coversfloatation = project.Project{
	ID:          "15219",
	Caption:     "Covers floatation",
	Customer:    "ACME",
	Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Fugit alias, totam corrupti eveniet nulla vero similique dignissimos. At quis officia omnis assumenda, quod dolore explicabo blanditiis, deleniti, nesciunt quibusdam nam.",
	Status:      project.Active,

	Engineers: []user.ID{},
	Estimate:  60 * day,

	Created:   time.Now().Add(-day * 3),
	Modified:  time.Now().Add(-day * 1),
	Completed: time.Time{},
}
