package project

import "time"

type Project struct {
	Title       string
	Customer    string //TODO: ref
	Pricing     Pricing
	Description string
	Status      Status
}

type Status string

const (
	Queued     Status = "Queued"
	InProgress        = "In Progress"
	Done              = "Done"
	Delivered         = "Delivered"
)

//TODO: is there a better name for this?
type Pricing struct {
	Hours float64
	Price float64
}

type Task struct {
	Title       string
	Description string
	Status      Status
}

type Resource struct {
	Name string
	Unit Unit
	PPU  float64 // price per unit
}

type Unit string

const (
	Litre = "l"
	Grams = "g"
	Hour  = "h"
	Piece = "unit"
)

type Expense struct {
	Worker   string //TODO: ref
	Date     time.Time
	Resource Resource
	Units    float64
	Price    float64
}
