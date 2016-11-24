package work

import "time"

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
	// MarkSubmitted marks the activities as submitted
	MarkSubmitted(activity []ActivityID) error
}

type ActivityID uint64

type Activity struct {
	ID   ActivityID
	Name string

	Start  time.Time
	Finish time.Time
}

func (activity *Activity) Duration() time.Duration {
	if activity.Finish.IsZero() {
		return time.Now().Sub(activity.Start)
	}
	return activity.Finish.Sub(activity.Start)
}
