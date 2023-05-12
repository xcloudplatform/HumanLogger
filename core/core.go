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

	screenshotsAttemptsStream := makeScreenshotsAttemptStream(&c.UserActivity)
	screenshotStream := makeScreenshotStream(screenshotsAttemptsStream)
	deduplicatedScreenshotStream := makeDeduplicatedScreenshotStream(screenshotStream)
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
			case isActive := <-userActivity.isActiveChan:
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

func makeDeduplicatedScreenshotStream(screenshotStream chan Screenshot) chan Screenshot {
	// Create a new channel to hold deduplicated screenshots
	deduplicatedStream := make(chan Screenshot)

	// Keep track of the last screenshot received for each display ID
	lastScreenshots := make(map[int]Screenshot)

	// Start a goroutine to read from the input stream and deduplicate the screenshots
	go func() {
		for screenshot := range screenshotStream {
			diffScreenshot := screenshot

			// Check if this screenshot is a duplicate of the last one we received for this display ID
			if lastScreenshot, ok := lastScreenshots[screenshot.DisplayID]; ok {

				same, diffScreenshotPtr, err := screenshot.Diff(&lastScreenshot)

				if err != nil || same {
					continue
				}
				diffScreenshot = *diffScreenshotPtr
			}

			// This is a new screenshot for this display ID, so send it on the deduplicated stream
			deduplicatedStream <- diffScreenshot

			// Update the last screenshot map
			lastScreenshots[screenshot.DisplayID] = screenshot
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
