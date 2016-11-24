package work

import "time"

type Summary struct {
	Start  time.Time
	Finish time.Time

	Activities []ActivityID
	Durations  map[string]time.Duration
}

func (summary *Summary) Include(activity Activity) error {
	if activity.Incomplete() {
		return ErrActivityIncomplete
	}

	if summary.Start.After(activity.Start) {
		summary.Start = activity.Start
	}
	if summary.Finish.Before(activity.Finish) {
		summary.Finish = activity.Finish
	}

	summary.Activities = append(summary.Activities, activity.ID)
	summary.Durations[activity.Name] += activity.Duration()

	return nil
}

func SummarizeActivities(activities []Activity) (*Summary, error) {
	summary := &Summary{}
	summary.Durations = make(map[string]time.Duration)

	for _, activity := range activities {
		if err := summary.Include(activity); err != nil {
			return nil, err
		}
	}

	return summary, nil
}
