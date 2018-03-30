package work

import (
	"errors"
	"sync"
	"time"

	"github.com/loov/timeclock/project"
)

var (
	ErrSheetDoesNotExist = errors.New("sheet does not exist")
)

var _ Sheets = &Database{}

type Database struct {
	mu sync.Mutex

	lastSheetID    SheetID
	lastActivityID ActivityID

	sheets     []*Sheet
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
	return ErrSheetDoesNotExist
}

func (db *Database) Delete(entryID SheetID) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	index := db.findSheetIndex(entryID)
	if index >= 0 {
		db.sheets = append(db.sheets[:index], db.sheets[index+1:]...)
		return nil
	}
	return ErrSheetDoesNotExist
}
