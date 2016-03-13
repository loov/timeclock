package db

import "github.com/loov/timeclock/dayreport"

func (db *DB) DayReports() dayreport.Reports {
	return &DayReports{}
}

type DayReports struct{}
