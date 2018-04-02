package work

import (
	"html/template"
	"time"
)

type Calendar struct {
	Days []CalendarDay
}

type CalendarDay struct {
	Day   time.Time
	Class []string
	Text  string
	Link  string
	Badge string
}

func NewCalendar() *Calendar {
	return &Calendar{}
}

func (calendar *Calendar) Render() template.HTML {
	return ""
}
