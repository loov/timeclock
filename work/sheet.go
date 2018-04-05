package work

import (
	"sort"
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

// Activities represents storage for work.Activities-s
type Activities interface {
	// WorkerSheet returns activities for a worker
	WorkerSheet(worker user.ID, start, end time.Time) (*Sheet, error)

	// Submit adds a new entry
	Submit(activities []Activity) error
	// Update updates existing activies with the appropriate ID-s
	Update(activities []Activity) error
	// Delete deletes activities
	Delete(activities []Activity) error
}

// ActivityID is an unique identifier for activity
type ActivityID int

type ActivityName string

// Activity represents one activity in a work.Sheet
type Activity struct {
	ID ActivityID

	// Time associated with this entry
	Time time.Time
	// Modified is last modified time for this activity
	Modified time.Time
	// Worker who created this entry
	Worker user.ID
	// Associated project
	Project project.ID

	// Name of the activity
	Name ActivityName
	// Duration is the time spent on this activity
	Duration time.Duration

	// Locked is a computed property showing it shouldn't be modified further
	Locked bool
}

// Sheet represents a day/week list of activities
type Sheet struct {
	// Dates associated
	Start, End time.Time
	// Latest modified time
	Modified time.Time

	// Worker contains a worker, when there is only one worker for all activities
	Worker user.ID

	// Project contains a project, when there is only one project for all activities
	Project project.ID

	// list of activities associated with this sheet
	Activities []Activity

	// total duration of the activities
	Duration time.Duration
}

func NewSheet(activities []Activity) *Sheet {
	sheet := &Sheet{}
	sheet.Activities = activities
	if len(activities) == 0 {
		return sheet
	}

	act := activities[0]
	sheet.Start, sheet.End = act.Time, act.Time
	sheet.Modified = act.Modified
	sheet.Worker = act.Worker
	sheet.Project = act.Project
	sheet.Duration = act.Duration

	for _, act := range activities[1:] {
		if act.Modified.After(sheet.Modified) {
			sheet.Modified = act.Modified
		}
		if act.Time.Before(sheet.Start) {
			sheet.Start = act.Time
		}
		if act.Time.After(sheet.End) {
			sheet.End = act.Time
		}
		if sheet.Worker != act.Worker {
			sheet.Worker = 0
		}
		if sheet.Project != act.Project {
			sheet.Project = 0
		}
		sheet.Duration += act.Duration
	}

	return sheet
}

func (sheet *Sheet) Before(b *Sheet) bool {
	if sheet.Start.Equal(b.Start) {
		return sheet.End.Before(b.End)
	}
	return sheet.Start.Before(b.Start)
}

func (sheet *Sheet) SummarizeByDay() []*Sheet {
	byDay := map[time.Time][]Activity{}
	for _, act := range sheet.Activities {
		day := Day(act.Time.UTC())
		byDay[day] = append(byDay[day], act)
	}

	list := []*Sheet{}
	for _, activities := range byDay {
		list = append(list, NewSheet(activities))
	}
	sort.Slice(list, func(i, k int) bool {
		return list[i].Before(list[k])
	})

	return list
}

func Day(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
