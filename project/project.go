package project

import (
	"errors"
	"time"

	"github.com/loov/timeclock/user"
)

var (
	ErrNotExist = errors.New("Project does not exist.")
)

type Status string

const (
	Inactive  Status = "Inactive"
	Active           = "Active"
	Delivered        = "Delivered"
)

type ID string

type Project struct {
	ID          ID
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
