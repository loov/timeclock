package work

import (
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

// Sheets represents storage for work.Sheet-s
type Sheets interface {
	// Overview returns information about day status
	Overview(from, to time.Time) ([]Overview, error)
	FullOverview(from, to time.Time) ([]FullOverview, error)

	// Submit adds a new entry
	Submit(sheet *Sheet) (SheetID, error)
	Update(sheet *Sheet) error
	Delete(sheetID SheetID) error
}

// SheetID is an unique identifier for an entry
type SheetID uint64

// Sheet represents a day/week submission by a worker
type Sheet struct {
	// ID is the unique id
	ID SheetID
	// Date associated
	Date time.Time
	// Worker who created this entry
	Worker user.ID

	// UpdatedAt is the last time this entry was modified
	UpdatedAt time.Time

	// list of activities associated with this sheet
	Activities []*Activity
	// total duration of the sheet
	Duration time.Duration

	// Locked means that the entry cannot be modified without special permissions
	Locked bool
}

// ActivityID is an unique identifier for activity
type ActivityID uint64

// Activity represents one activity in a work.Sheet
type Activity struct {
	SheetID SheetID
	ID      ActivityID

	// Date associated
	Date time.Time
	// Worker who created this entry
	Worker user.ID
	// Associated project
	Project project.ID

	// Name of the activity
	Name string
	// Duration is the time spent on this activity
	Duration time.Duration
}

type Overview struct {
	Date   time.Time
	Total  time.Duration
	Locked bool
}

type FullOverview struct {
	Overview

	Sheets []SheetID

	ByWorker   map[user.ID]time.Duration
	ByActivity map[string]time.Duration
	ByProject  map[project.ID]time.Duration
}
