package work

import (
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

	// Worker contains a worker, when there is only one worker for all activities
	Worker user.ID

	// Project contains a project, when there is only one project for all activities
	Project project.ID

	// list of activities associated with this sheet
	Activities []Activity

	// total duration of the activities
	Duration time.Duration
}

type Summary struct {
	Worker   map[user.ID]time.Duration
	Activity map[ActivityName]time.Duration
	Project  map[project.ID]time.Duration
}
