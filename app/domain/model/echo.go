package model

import (
	"time"

	"github.com/google/uuid"
)

type Echo struct {
	ID        uuid.UUID
	Message   string
	Timestamp time.Time
}

func NewEcho(message string, timestamp time.Time) Echo {
	return Echo{
		ID:        uuid.New(),
		Message:   message,
		Timestamp: timestamp,
	}
}
