package main

import (
	"fmt"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
)

type UserActivity struct {
	isActive      bool
	isActiveChan  chan bool
	eventsChan    chan hook.Event
	lastEventTime time.Time
	mux           sync.Mutex
}

func (ua *UserActivity) start() {
	ua.isActiveChan = make(chan bool, 11)
	ua.eventsChan = make(chan hook.Event, 11)

	go func() {
		evChan := hook.Start()
		defer hook.End()

		// start a new ticker to set inactive in future
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case ev := <-evChan:
				//fmt.Println("hook: ", ev)
				ua.eventsChan <- ev
				ua.setActive(true)
				ticker.Stop()
				ticker = time.NewTicker(1 * time.Second)
			case <-ticker.C:
				ua.setActive(false)
			}
		}
	}()

}

func (ua *UserActivity) setActive(active bool) {
	ua.mux.Lock()
	defer ua.mux.Unlock()

	if active != ua.isActive {
		ua.isActive = active
		ua.isActiveChan <- active
	}
	if active {
		ua.lastEventTime = time.Now()
	}
}

func (ua *UserActivity) isActiveNow() bool {
	ua.mux.Lock()
	defer ua.mux.Unlock()

	return ua.isActive
}

func main() {
	userActivity := &UserActivity{}

	userActivity.start()

	isScreenshotTakenChan := takeScreenshots(userActivity)
	getWindowTitles(isScreenshotTakenChan)
	//select {}
	evChan := filterEvents(userActivity.eventsChan, userActivity)
	for ev := range evChan {
		fmt.Println("filtered: ", ev)

	}
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

func takeScreenshots(userActivity *UserActivity) chan bool {
	isScreenshotTakenChan := make(chan bool, 2000)

	go func() {
		for {
			isActive := userActivity.isActiveNow()

			interval := 5 * time.Second
			if isActive {
				interval = 300 * time.Millisecond
			}
			fmt.Println("Sleeping for", interval)

			select {
			case <-userActivity.isActiveChan:
				// Take a screenshot immediately when isActiveChan is changed
				fmt.Println("TAKING SCREENSHOT! (activeChan changed)")
				//fixme compare screenshot
				isScreenshotTakenChan <- true

			case <-time.After(interval):
				// Take a screenshot after sleeping for the appropriate interval
				fmt.Println("TAKING SCREENSHOT! (after sleeping)")
				//fixme compare screenshot
				isScreenshotTakenChan <- true
			}
		}
	}()

	return isScreenshotTakenChan
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
