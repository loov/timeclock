package work

import (
	"errors"
	"sync"
	"time"

	"github.com/loov/timeclock/project"
)

var (
	ErrEntryDoesNotExist = errors.New("entry does not exist")
)

var _ Sheet = &Database{}

type Database struct {
	mu sync.Mutex

	lastEntryID    EntryID
	lastActivityID ActivityID

	entries    []*Entry
	activities []*Activity
	projects   []*project.Project
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) DefaultActivities() ([]string, error) {
	return []string{"Plumbing", "Welding", "Construction"}, nil
}

func (db *Database) Overview(from, to time.Time) ([]Overview, error) {
	return nil, nil
}
func (db *Database) FullOverview(from, to time.Time) ([]FullOverview, error) {
	return nil, nil
}

func (db *Database) Submit(entry *Entry) (EntryID, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	entry.UpdatedAt = time.Now()

	db.lastEntryID++
	entry.ID = db.lastEntryID

	for _, activity := range entry.Activities {
		db.lastActivityID++
		activity.ID = db.lastActivityID
		activity.Entry = entry.ID
		activity.Date = entry.Date
		activity.Worker = entry.Worker

		db.activities = append(db.activities, activity)
	}

	db.entries = append(db.entries, entry)

	return entry.ID, nil
}

func (db *Database) Update(entry *Entry) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, en := range db.entries {
		if en.ID == entry.ID {
			db.entries[i] = entry
			return nil
		}
	}

	return ErrEntryDoesNotExist
}

func (db *Database) Delete(entryID EntryID) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, en := range db.entries {
		if en.ID == entryID {
			db.entries = append(db.entries[:i], db.entries[i+1:]...)
			return nil
		}
	}

	return ErrEntryDoesNotExist
}
