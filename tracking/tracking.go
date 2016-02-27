package tracking

import (
	"time"

	"github.com/loov/workclock/user"
)

type ActivityID int64

type Activity struct {
	Project   project.ID
	Worker    user.ID
	Activity  string
	Start     time.Time
	Finish    time.Time
	Processed bool
}

type Service interface {
	SelectProject(user user.ID, project project.ID)
	SelectActivity(user user.ID, project project.ID, activity string)
	FinishActivity(user user.ID, project project.ID, activity string)
}

type Activities interface {
	Unprocessed() []Activity
	UnprocessedByUser(user user.ID) []Activity

	MarkProcessed(activity []ActivityID)
}
