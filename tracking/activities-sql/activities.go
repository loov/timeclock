package activities

import (
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/tracking"
	"github.com/loov/timeclock/user"
)

var (
	_ tracking.Activities = &Activities{}
	_ tracking.Tracker    = &Activities{}
)

type Activities struct{}

func New(connection string) (*Activities, error) {
	return &Activities{}, nil
}

func (_ Activities) SelectActivity(user user.ID, project project.ID, activity string) error {
	return nil
}
func (_ Activities) FinishActivity(user user.ID, project project.ID, activity string) error {
	return nil
}
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
