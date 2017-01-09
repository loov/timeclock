package work

import "errors"

var (
	ErrActivityIncomplete = errors.New("activity is incomplete.")
	ErrNoCurrentActivity  = errors.New("no current activity")
)

type Activities interface {
	// Names returns list of available activities
	Names() ([]string, error)

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
