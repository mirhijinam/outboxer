package model

import "time"

type Event struct {
	ID          int
	Payload     string
	ReservedFor time.Duration
}
