package work

import (
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

type ActivityID uint64

type Activity struct {
	// ID is an unique identifier for this entry
	ID ActivityID
	// User who is performing the activity
	Worker user.ID
	// Project where this activity is performed
	Project project.ID

	// Name of the this activity
	Name string
	// Start is the start time of the activity in UTC
	Start time.Time
	// Finish is the finishing time of the activity in UTC
	Finish time.Time
	// Submitted marks the activity as sent for review
	Submitted bool
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
