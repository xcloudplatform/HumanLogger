package core

import "time"

type Screenshot struct {
	ID        string
	UserID    string
	Data      []byte
	Timestamp time.Time
}
