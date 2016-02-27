package tracking

import (
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

type Activity struct {
	Project   project.ID
	Worker    user.ID
	Activity  string
	Start     time.Time
	Finish    time.Time
	Processed bool
}

// TODO: attach to user
var ActivityNames = []string{
	"welding",
	"plumbing",
	"plateworks",
	"other",
}
