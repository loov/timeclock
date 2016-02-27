package activities

import (
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
func (_ Activities) Unprocessed() ([]tracking.Activity, error) {
	return nil, nil
}
func (_ Activities) UnprocessedByUser(user user.ID) ([]tracking.Activity, error) {
	return nil, nil
}
func (_ Activities) MarkProcessed(activity []tracking.Activity) error {
	return nil
}
