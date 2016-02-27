package project

import (
	"time"

	"github.com/loov/timeclock/user"
)

type Status string

const (
	Inactive  Status = "Inactive"
	Active           = "Active"
	Delivered        = "Delivered"
)

type Project struct {
	ID          string
	Caption     string
	Customer    string
	Description string
	Status      Status

	Engineers []user.ID
	Estimate  time.Duration

	Created   time.Time
	Modified  time.Time
	Completed time.Time
}

type Projects interface {
	List() []Project
}
