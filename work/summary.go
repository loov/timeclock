package work

import (
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

type SummaryID uint64

type Summary struct {
	ID      SummaryID
	Worker  user.ID
	Project project.ID

	Start  time.Time
	Finish time.Time

	Activities []ActivityID
	Durations  map[string]time.Duration
}

func NewSummary() *Summary {
	summary := &Summary{}
	summary.Durations = make(map[string]time.Duration)
	return summary
}

func (summary *Summary) Include(activity Activity) error {
	if activity.Incomplete() {
		return ErrActivityIncomplete
	}

	if summary.Start.IsZero() || summary.Start.After(activity.Start) {
		summary.Start = activity.Start
	}
	if summary.Finish.IsZero() || summary.Finish.Before(activity.Finish) {
		summary.Finish = activity.Finish
	}

	summary.Activities = append(summary.Activities, activity.ID)
	summary.Durations[activity.Name] += activity.Duration()

	return nil
}

func SummarizeActivities(activities []Activity) (*Summary, error) {
	summary := NewSummary()
	for _, activity := range activities {
		if err := summary.Include(activity); err != nil {
			return nil, err
		}
	}
	return summary, nil
}

func SummarizeActivitiesByDay(activities []Activity, loc *time.Location) (map[time.Time]*Summary, error) {
	summaries := map[time.Time]*Summary{}
	for _, activity := range activities {
		day := activity.Start.In(loc).Truncate(24 * time.Hour)

		summary, ok := summaries[day]
		if !ok {
			summary = NewSummary()
			summaries[day] = summary
		}

		if err := summary.Include(activity); err != nil {
			return nil, err
		}
	}

	return summaries, nil
}
