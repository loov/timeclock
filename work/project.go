package work

import (
	"sync"
	"time"
)

type Project struct {
	mu         sync.Mutex
	activities []string
	jobs       []Job
	days       []Day
}

func NewProject() *Project {
	project := &Project{}
	project.activities = []string{"Plumbing", "Welding", "Construction"}
	return project
}

func (project *Project) Activities() []string {
	project.mu.Lock()
	defer project.mu.Unlock()
	return append([]string{}, project.activities...)
}

func (project *Project) Jobs() []Job {
	project.mu.Lock()
	defer project.mu.Unlock()

	return append([]Job{}, project.jobs...)
}

func (project *Project) Days() []Day {
	project.mu.Lock()
	defer project.mu.Unlock()

	return append([]Day{}, project.days...)
}

func (project *Project) SelectActivity(activity string) {
	project.mu.Lock()
	defer project.mu.Unlock()

	now := time.Now()
	if len(project.jobs) > 0 {
		last := &project.jobs[len(project.jobs)-1]
		if last.Finish.IsZero() {
			last.Finish = now
		}
	}

	if activity != "" {
		project.jobs = append(project.jobs, Job{
			Activity: activity,
			Start:    time.Now(),
		})
	}
}

func (project *Project) SubmitDay() {
	project.mu.Lock()
	defer project.mu.Unlock()

	durations := map[string]time.Duration{}
	for _, job := range project.jobs {
		durations[job.Activity] += job.Duration()
	}

	day := Day{
		Submitted:  time.Now(),
		Activities: durations,
	}

	project.days = append(project.days, day)
	project.jobs = nil
}

func (project *Project) Summary() map[string]time.Duration {
	project.mu.Lock()
	defer project.mu.Unlock()

	durations := map[string]time.Duration{}
	for _, job := range project.jobs {
		durations[job.Activity] += job.Duration()
	}
	return durations
}

func (project *Project) CurrentActivity() string {
	project.mu.Lock()
	defer project.mu.Unlock()

	if len(project.jobs) == 0 {
		return ""
	}

	last := &project.jobs[len(project.jobs)-1]
	if last.Finish.IsZero() {
		return last.Activity
	}

	return ""
}
