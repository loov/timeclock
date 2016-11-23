package work

import (
	"sync"
	"time"
)

type Model struct {
	mu         sync.Mutex
	activities []string
	jobs       []Job
	days       []Day
}

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

func (model *Model) Activities() []string {
	model.mu.Lock()
	defer model.mu.Unlock()
	return append([]string{}, model.activities...)
}

func (model *Model) Jobs() []Job {
	model.mu.Lock()
	defer model.mu.Unlock()

	return append([]Job{}, model.jobs...)
}

func (model *Model) Days() []Day {
	model.mu.Lock()
	defer model.mu.Unlock()

	return append([]Day{}, model.days...)
}

func (model *Model) SelectActivity(activity string) {
	model.mu.Lock()
	defer model.mu.Unlock()

	now := time.Now()
	if len(model.jobs) > 0 {
		last := &model.jobs[len(model.jobs)-1]
		if last.Finish.IsZero() {
			last.Finish = now
		}
	}

	if activity != "" {
		model.jobs = append(model.jobs, Job{
			Activity: activity,
			Start:    time.Now(),
		})
	}
}

func (model *Model) SubmitDay() {
	model.mu.Lock()
	defer model.mu.Unlock()

	durations := map[string]time.Duration{}
	for _, job := range model.jobs {
		durations[job.Activity] += job.Duration()
	}

	day := Day{
		Submitted:  time.Now(),
		Activities: durations,
	}

	model.days = append(model.days, day)
	model.jobs = nil
}

func (model *Model) Summary() map[string]time.Duration {
	model.mu.Lock()
	defer model.mu.Unlock()

	durations := map[string]time.Duration{}
	for _, job := range model.jobs {
		durations[job.Activity] += job.Duration()
	}
	return durations
}

func (model *Model) CurrentActivity() string {
	model.mu.Lock()
	defer model.mu.Unlock()

	if len(model.jobs) == 0 {
		return ""
	}

	last := &model.jobs[len(model.jobs)-1]
	if last.Finish.IsZero() {
		return last.Activity
	}

	return ""
}
