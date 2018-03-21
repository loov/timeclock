package work

import (
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

// Sheet represents a work-sheet
type Sheet interface {
	// Overview returns information about day status
	Overview(from, to time.Time) ([]Overview, error)
	FullOverview(from, to time.Time) ([]FullOverview, error)

	// Submit adds a new entry
	Submit(entry *Entry) (EntryID, error)
	Update(entry *Entry) error
	Delete(entry EntryID) error
}

// EntryID is an unique identifier for an entry
type EntryID uint64

// Entry represents a day/week submission by a worker
type Entry struct {
	// ID is the unique id
	ID EntryID
	// Date associated
	Date time.Time
	// Worker who created this entry
	Worker user.ID

	// UpdatedAt is the last time this entry was modified
	UpdatedAt time.Time

	// list of activities associated with this Entry
	Activities []*Activity
	// total duration of the entry
	Duration time.Duration

	// Locked means that the entry cannot be modified without special permissions
	Locked bool
}

// ActivityID is an unique identifier for activity
type ActivityID uint64

// Activity represents one activity in a work.Entry
type Activity struct {
	Entry EntryID
	ID    ActivityID

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

	Entries []EntryID

	ByWorker   map[user.ID]time.Duration
	ByActivity map[string]time.Duration
	ByProject  map[project.ID]time.Duration
}
