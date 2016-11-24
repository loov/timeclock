package work

import (
	"errors"
	"time"
)

var (
	ErrActivityIncomplete = errors.New("activity is incomplete.")
	ErrNoCurrentActivity  = errors.New("no current activity")
)

type Activities interface {
	// DefaultNames returns the default list of activities
	DefaultNames() ([]string, error)

	// Current returns the current activity
	Current() (Activity, error)
	// Start starts a new activity and finishes the previous and starts a new activity
	Start(activity string) error
	// Finish finishes the current activity
	Finish() error

	// Pending returns the list of activities that have not been marked as reported
	Pending() ([]Activity, error)

	// Report marks the summary as submitted
	Report(summary *Summary) error
	// Reports returns the list of submitted reports
	Reports() ([]*Summary, error)
}

type ActivityID uint64

type Activity struct {
	ID   ActivityID
	Name string

	Start  time.Time
	Finish time.Time
}

func (activity *Activity) Incomplete() bool {
	return activity.Start.IsZero() || activity.Finish.IsZero()
}

func (activity *Activity) Duration() time.Duration {
	if activity.Finish.IsZero() {
		return time.Now().Sub(activity.Start)
	}
	return activity.Finish.Sub(activity.Start)
}
