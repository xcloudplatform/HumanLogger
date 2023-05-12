package core

import (
	hook "github.com/robotn/gohook"
	"time"
)

type UIEvent struct {
	Timestamp time.Time `json:"timestamp"`

	hook.Event
}

func NewUIEvent(event hook.Event) UIEvent {
	return UIEvent{
		Timestamp: time.Now(),
		Event:     event,
	}
}
