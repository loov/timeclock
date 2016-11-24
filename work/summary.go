package work

import (
	"errors"
	"time"
)

type Summary struct {
	Start  time.Time
	Finish time.Time

	Activities []ActivityID
	Durations  map[string]time.Duration
}

func SummarizeActivities(activities []Activity) (*Summary, error) {
	summary := &Summary{}
	summary.Durations = make(map[string]time.Duration)

	for _, activity := range activities {
		// don't count unfinished activities
		if activity.Start.IsZero() || activity.Finish.IsZero() {
			return nil, errors.New("activity is incomplete")
		}

		if summary.Start.After(activity.Start) {
			summary.Start = activity.Start
		}
		if summary.Finish.Before(activity.Finish) {
			summary.Finish = activity.Finish
		}

		summary.Activities = append(summary.Activities, activity.ID)
		summary.Durations[activity.Name] += activity.Duration()
	}

	return summary, nil
}
