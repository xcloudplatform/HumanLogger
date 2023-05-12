package core

import (
	"time"

	hook "github.com/robotn/gohook"
)

type Core struct {
	UserActivity     UserActivity
	ScreenshotStream chan Screenshot
	UIEventStream    chan UIEvent
}

func NewCore() *Core {
	return &Core{
		UserActivity:     UserActivity{},
		ScreenshotStream: make(chan Screenshot),
		UIEventStream:    make(chan UIEvent),
	}
}

func (c *Core) Start() error {
	c.UserActivity.start()

	isScreenshotTakenChan := takeScreenshotsAttemptStream(&c.UserActivity)
	getWindowTitles(isScreenshotTakenChan)

	evChan := filterEvents(c.UserActivity.eventsChan, &c.UserActivity)
	for ev := range evChan {
		uiEv := NewUIEvent(ev)
		c.UIEventStream <- uiEv
	}

	return nil
}

func (c *Core) Stop() {
	// TODO: Implement the Stop method.
}

func filterEvents(events chan hook.Event, userActivity *UserActivity) chan hook.Event {
	filteredEvents := make(chan hook.Event, 20)
	var prevEvent hook.Event
	var firstMouseMoveEvent hook.Event

	go func() {
		for {
			select {
			case ev := <-events:
				if ev.Kind == hook.MouseMove {
					if prevEvent.Kind != hook.MouseMove {
						// First MouseMove event in the sequence
						firstMouseMoveEvent = ev
						filteredEvents <- ev
					}
				} else {
					if prevEvent.Kind == hook.MouseMove {
						// Last MouseMove event in the sequence
						filteredEvents <- firstMouseMoveEvent
					}
					filteredEvents <- ev
				}
				prevEvent = ev
			case isActive := <-userActivity.isActiveChan:
				if !isActive {
					if prevEvent.Kind == hook.MouseMove {
						filteredEvents <- prevEvent

					}
				}

			}

		}
	}()

	return filteredEvents
}

func takeScreenshotsAttemptStream(userActivity *UserActivity) chan bool {
	stream := make(chan bool, 2000)

	go func() {
		for {
			isActive := userActivity.isActiveNow()

			interval := 5 * time.Second
			if isActive {
				interval = 300 * time.Millisecond
			}
			//fmt.Println("Sleeping for", interval)

			select {
			case <-userActivity.isActiveChan:
				// Take a screenshot immediately when isActiveChan is changed
				//fmt.Println("TAKING SCREENSHOT! (activeChan changed)")
				stream <- true

			case <-time.After(interval):
				// Take a screenshot after sleeping for the appropriate interval
				//fmt.Println("TAKING SCREENSHOT! (after sleeping)")
				stream <- true
			}
		}
	}()

	return stream
}

func getWindowTitles(isScreenshotTakenChan chan bool) {
	go func() {
		for {
			isScreenshotTaken := <-isScreenshotTakenChan
			if isScreenshotTaken {
				//title := robotgo.GetTitle()
				//fmt.Println("title@@@ ", title)

			}

		}
	}()
}
