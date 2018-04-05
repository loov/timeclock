package work

import (
	"errors"
	"sync"
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

var (
	ErrActivityDoesNotExist = errors.New("activity does not exist")
)

var _ Activities = &Database{}

type Database struct {
	mu sync.Mutex

	lastActivityID ActivityID

	activities map[ActivityID]Activity
	projects   []*project.Project
}

func NewDatabase() *Database {
	db := &Database{}
	db.activities = make(map[ActivityID]Activity)
	return db
}

func (db *Database) DefaultActivities() ([]string, error) {
	return []string{"Plumbing", "Welding", "Construction"}, nil
}

func (db *Database) createSheet(worker user.ID, project project.ID, start, end time.Time) *Sheet {
	activities := []Activity{}
	for _, act := range db.activities {
		if worker != 0 && worker != act.Worker {
			continue
		}
		if project != 0 && project != act.Project {
			continue
		}
		if !start.IsZero() && act.Time.Before(start) {
			continue
		}
		if !end.IsZero() && act.Time.After(end) {
			continue
		}

		activities = append(activities, act)
	}

	sheet := NewSheet(activities)
	sheet.Start, sheet.End = start, end

	return sheet
}

// WorkerSheet returns activities for a worker
func (db *Database) WorkerSheet(worker user.ID, start, end time.Time) (*Sheet, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.createSheet(worker, 0, start, end), nil
}

// Submit adds a new entry
func (db *Database) Submit(activities []Activity) error {
	modified := time.Now()

	db.mu.Lock()
	defer db.mu.Unlock()

	for i := range activities {
		act := &activities[i]
		db.lastActivityID++
		act.ID = db.lastActivityID
		act.Modified = modified
		db.activities[act.ID] = *act
	}

	return nil
}

// Update updates existing activies with the appropriate ID-s
func (db *Database) Update(activities []Activity) error {
	modified := time.Now()

	db.mu.Lock()
	defer db.mu.Unlock()

	for _, act := range activities {
		_, ok := db.activities[act.ID]
		if !ok {
			return ErrActivityDoesNotExist
		}
	}

	for i := range activities {
		act := &activities[i]
		act.Modified = modified
		db.activities[act.ID] = *act
	}

	return nil
}

// Delete deletes activities
func (db *Database) Delete(activities []Activity) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, act := range activities {
		_, ok := db.activities[act.ID]
		if !ok {
			return ErrActivityDoesNotExist
		}
	}

	for i := range activities {
		act := &activities[i]
		delete(db.activities, act.ID)
		act.ID = 0
	}

	return nil
}
