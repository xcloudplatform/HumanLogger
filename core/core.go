package core

import (
	"fmt"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
)

type Core struct {
	UserActivity     UserActivity
	ScreenshotStream chan Screenshot
	UIEventStream    chan UIEvent

	lastScreenshots map[int]Screenshot
	mux             sync.Mutex
}

func NewCore() *Core {
	return &Core{
		lastScreenshots: make(map[int]Screenshot),

		UserActivity:     UserActivity{},
		ScreenshotStream: make(chan Screenshot),
		UIEventStream:    make(chan UIEvent),
	}
}

func (c *Core) Start() error {
	c.UserActivity.start()

	screenshotsAttemptsStream := makeScreenshotsAttemptStream(&c.UserActivity)
	screenshotStream := makeScreenshotStream(screenshotsAttemptsStream)
	deduplicatedScreenshotStream := c.makeDeduplicatedScreenshotStream(screenshotStream)
	getWindowTitles(screenshotsAttemptsStream)

	go func() {
		for scr := range deduplicatedScreenshotStream {
			c.ScreenshotStream <- scr
		}

	}()

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
			case isActive := <-userActivity.IsActiveStream:
				if !isActive {
					if prevEvent.Kind == hook.MouseMove {
						filteredEvents <- prevEvent

					}
				}

			}

		}
		close(filteredEvents)
	}()

	return filteredEvents
}

func makeScreenshotsAttemptStream(userActivity *UserActivity) chan bool {
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
			case <-userActivity.IsActiveStream:
				// Take a screenshot immediately when IsActiveStream is changed
				//fmt.Println("TAKING SCREENSHOT! (activeChan changed)")
				stream <- true

			case <-time.After(interval):
				// Take a screenshot after sleeping for the appropriate interval
				//fmt.Println("TAKING SCREENSHOT! (after sleeping)")
				stream <- true
			}
		}
		close(stream)
	}()

	return stream
}

func makeScreenshotStream(screenshotsAttemptsStream chan bool) chan Screenshot {
	stream := make(chan Screenshot)
	go func() {
		for {
			isScreenshotTaken := <-screenshotsAttemptsStream
			if isScreenshotTaken {

				if screenshots, err := CaptureScreenshots(); err == nil {

					for _, scr := range screenshots {
						stream <- scr
					}
				}
			}

		}
		close(stream)
	}()
	return stream
}

func (c *Core) ResetLastScreenshots() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.lastScreenshots = make(map[int]Screenshot)
}

func (c *Core) makeDeduplicatedScreenshotStream(screenshotStream chan Screenshot) chan Screenshot {
	// Create a new channel to hold deduplicated screenshots
	deduplicatedStream := make(chan Screenshot)

	// Start a goroutine to read from the input stream and deduplicate the screenshots
	go func() {
		for screenshot := range screenshotStream {
			var diffScreenshot Screenshot
			c.mux.Lock()
			// Check if this screenshot is a duplicate of the last one we received for this display ID
			if lastScreenshot, ok := c.lastScreenshots[screenshot.DisplayID]; ok {

				same, diffScreenshotPtr, err := screenshot.Diff(&lastScreenshot)

				if err != nil || same {
					c.mux.Unlock()
					continue
				}
				diffScreenshot = *diffScreenshotPtr
				timestamp := time.Now().Format("15:04:05.000")

				fmt.Printf("%s: frame\n", timestamp)

			} else {
				diffScreenshot = screenshot
				timestamp := time.Now().Format("15:04:05.000")
				fmt.Printf("%s: new keyframe\n", timestamp)
			}

			// This is a new screenshot for this display ID, so send it on the deduplicated stream
			deduplicatedStream <- diffScreenshot

			// Update the last screenshot map

			c.lastScreenshots[screenshot.DisplayID] = screenshot
			c.mux.Unlock()
		}

		// Close the deduplicated stream when the input stream is closed
		close(deduplicatedStream)
	}()

	return deduplicatedStream
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
