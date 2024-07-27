package model

import "time"

type Message struct {
	ID        int
	Content   string
	CreatedAt time.Time
}
