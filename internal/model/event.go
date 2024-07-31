package model

import "time"

type Event struct {
	ID            int
	Payload       string
	ReservedUntil time.Time
}
