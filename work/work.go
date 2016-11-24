package work

import "time"

type Job struct {
	Activity string
	Start    time.Time
	Finish   time.Time
}

func (job *Job) Duration() time.Duration {
	if job.Finish.IsZero() {
		return time.Now().Sub(job.Start)
	}
	return job.Finish.Sub(job.Start)
}

type Day struct {
	Submitted  time.Time
	Activities map[string]time.Duration
}
