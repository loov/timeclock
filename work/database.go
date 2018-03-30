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
	sheet := &Sheet{}
	sheet.Start, sheet.End = start, end
	sheet.Worker, sheet.Project = worker, project

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

		sheet.Activities = append(sheet.Activities, act)
	}

	for i := range sheet.Activities {
		act := &sheet.Activities[i]
		sheet.Duration += act.Duration
	}

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

/*
func (db *Database) findSheetIndex(id SheetID) int {
	for i, sheet := range db.sheets {
		if sheet.ID == id {
			return i
		}
	}
	return -1
}

func (db *Database) Submit(sheet *Sheet) (SheetID, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	sheet.UpdatedAt = time.Now()

	db.lastSheetID++
	sheet.ID = db.lastSheetID

	for _, activity := range sheet.Activities {
		db.lastActivityID++
		activity.ID = db.lastActivityID
		activity.SheetID = sheet.ID
		activity.Date = sheet.Date
		activity.Worker = sheet.Worker

		db.activities = append(db.activities, activity)
	}

	db.sheets = append(db.sheets, sheet)

	return sheet.ID, nil
}

func (db *Database) Update(sheet *Sheet) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	index := db.findSheetIndex(sheet.ID)
	if index >= 0 {
		db.sheets[index] = sheet
		return nil
	}
	return ErrActivityDoesNotExist
}

func (db *Database) Delete(entryID SheetID) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	index := db.findSheetIndex(entryID)
	if index >= 0 {
		db.sheets = append(db.sheets[:index], db.sheets[index+1:]...)
		return nil
	}
	return ErrActivityDoesNotExist
}
*/
