package work

import (
	"math/rand"
	"sync"
	"time"
)

type ProjectID uint64

type Project struct {
	ID   ProjectID
	Name string

	ShortDescription string

	mu           sync.Mutex
	defaultNames []string
	lastID       ActivityID
	pending      []Activity
	submitted    []Activity
	reports      []*Summary
}

func NewProject(name, desc string) *Project {
	project := &Project{}
	project.ID = ProjectID(rand.Int())
	project.Name = name
	project.ShortDescription = desc
	project.defaultNames = []string{"Plumbing", "Welding", "Construction"}
	return project
}

func (project *Project) DefaultNames() ([]string, error) {
	project.mu.Lock()
	defer project.mu.Unlock()

	return append([]string{}, project.defaultNames...), nil
}

func (project *Project) Current() (Activity, error) {
	project.mu.Lock()
	defer project.mu.Unlock()

	if len(project.pending) > 0 {
		last := project.pending[len(project.pending)-1]
		if last.Incomplete() {
			return last, nil
		}
	}

	return Activity{}, ErrNoCurrentActivity
}

func (project *Project) _finishLast(now time.Time) {
	if len(project.pending) > 0 {
		last := &project.pending[len(project.pending)-1]
		if last.Finish.IsZero() {
			last.Finish = now
		}
	}
}

func (project *Project) Start(activity string) error {
	project.mu.Lock()
	defer project.mu.Unlock()

	now := time.Now()
	project._finishLast(now)

	if activity != "" {
		project.lastID++
		project.pending = append(project.pending, Activity{
			ID:    project.lastID,
			Name:  activity,
			Start: now,
		})
	}

	return nil
}

func (project *Project) Finish() error {
	project.mu.Lock()
	defer project.mu.Unlock()

	project._finishLast(time.Now())
	return nil
}

func (project *Project) Pending() ([]Activity, error) {
	project.mu.Lock()
	defer project.mu.Unlock()

	return append([]Activity{}, project.pending...), nil
}

func (project *Project) Report(summary *Summary) error {
	project.mu.Lock()
	defer project.mu.Unlock()

	//TODO: check duplicate submissions

	markSubmitted := map[ActivityID]struct{}{}
	for _, id := range summary.Activities {
		markSubmitted[id] = struct{}{}
	}

	for i := 0; i < len(project.pending); {
		act := project.pending[i]
		if _, ok := markSubmitted[act.ID]; ok {
			project.submitted = append(project.submitted, act)
			project.pending = append(project.pending[:i], project.pending[i+1:]...)
		} else {
			i++
		}
	}

	project.reports = append(project.reports, summary)

	return nil
}

func (project *Project) Reports() ([]*Summary, error) {
	project.mu.Lock()
	defer project.mu.Unlock()

	return append([]*Summary{}, project.reports...), nil
}
