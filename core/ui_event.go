package core

import "time"

type UIEvent struct {
	ID        string
	UserID    string
	Type      string
	Target    string
	Timestamp time.Time
}
