package db

import (
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/tracking"
	"github.com/loov/timeclock/user"
)

type Tracker struct{}

func (db *DB) Tracker() tracking.Tracker {
	return &Tracker{}
}

func (_ Tracker) SelectActivity(user user.ID, project project.ID, activity string) error {
	return nil
}

func (_ Tracker) FinishActivity(user user.ID, project project.ID, activity string) error {
	return nil
}

type Activities struct{}

func (db *DB) Activities() tracking.Activities { return &Activities{} }

func (_ Activities) Unreported() ([]tracking.Activity, error) {
	return []tracking.Activity{
		{
			"15324", 0, "welding",
			time.Now().Add(-15341 * time.Second),
			time.Now().Add(-14741 * time.Second),
			false,
		},
	}, nil
}

func (_ Activities) UnreportedByUser(user user.ID) ([]tracking.Activity, error) {
	return []tracking.Activity{
		{
			"15324", user, "welding",
			time.Now().Add(-15341 * time.Second),
			time.Now().Add(-14741 * time.Second),
			false,
		},
	}, nil
}
func (_ Activities) MarkProcessed(activity []tracking.Activity) error {
	return nil
}
