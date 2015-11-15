package project

import (
	"sort"
	"time"
)

type Activities []*Activity
type Materials []*Material

type Info struct {
	Project    *Project
	Activities Activities
	Materials  Materials
}

func (info *Info) Sort() {
	info.Activities.Sort()
	info.Materials.Sort()
}

type DayInfo struct {
	Date       time.Time
	Activities Activities
	Materials  Materials
}

func (info *Info) GroupByDay() []DayInfo {
	info.Sort()

	days := []DayInfo{}
	as, ms := info.Activities, info.Materials

	ai, mi := 0, 0
	for ai < len(as) || mi < len(ms) {
		day := DayInfo{}

		var next *time.Time
		if ai < len(as) && (next == nil || as[ai].Start.Before(*next)) {
			next = &as[ai].Start
		}
		if mi < len(ms) && (next == nil || ms[mi].Date.Before(*next)) {
			next = &ms[mi].Date
		}
		y, m, d := next.Date()

		day.Date = time.Date(y, m, d, 0, 0, 0, 0, next.Location())

		endsAt := day.Date.Add(24 * time.Hour)
		for ai < len(as) && as[ai].Start.Before(endsAt) {
			day.Activities = append(day.Activities, as[ai])
			ai++
		}
		for mi < len(ms) && ms[mi].Date.Before(endsAt) {
			day.Materials = append(day.Materials, ms[mi])
			mi++
		}

		days = append(days, day)
	}

	return days
}

func (xs Activities) Sort()              { sort.Sort(xs) }
func (xs Activities) Len() int           { return len(xs) }
func (xs Activities) Swap(i, j int)      { xs[i], xs[j] = xs[j], xs[i] }
func (xs Activities) Less(i, j int) bool { return xs[i].Date.Before(xs[j].Date) }

func (xs Materials) Sort()              { sort.Sort(xs) }
func (xs Materials) Len() int           { return len(xs) }
func (xs Materials) Swap(i, j int)      { xs[i], xs[j] = xs[j], xs[i] }
func (xs Materials) Less(i, j int) bool { return xs[i].Date.Before(xs[j].Date) }
