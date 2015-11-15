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

type Event struct {
	Worker string
	Date   time.Time
}

func (ev Event) Info() Event { return ev }

type Activity struct {
	Event
	Name    string
	Start   time.Time
	Finish  time.Time
	Comment string
}

func (a *Activity) Duration() time.Duration {
	return a.Finish.Sub(a.Start)
}

type Material struct {
	Event
	Resource Resource
	Amount   float64
}

func (m *Material) Total() float64 {
	return m.Resource.PPU * m.Amount
}
