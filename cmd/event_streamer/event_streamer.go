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
	lastEventTime time.Time
	mux           sync.Mutex
}

func (ua *UserActivity) start() {
	ua.isActiveChan = make(chan bool, 100)
	evChan := hook.Start()
	defer hook.End()

	// start a new ticker to set inactive in future
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case ev := <-evChan:
			fmt.Println("hook: ", ev)
			ua.setActive(true)
			ticker.Stop()
			ticker = time.NewTicker(1 * time.Second)
		case <-ticker.C:
			ua.setActive(false)
		}
	}
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

	go userActivity.start()

	isScreenshotTakenChan := takeScreenshots(userActivity)
	getWindowTitles(isScreenshotTakenChan)
	select {}
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
			time.Sleep(interval)

			fmt.Println("TAKING SCREENSHOT!")
			//fixme compare screenshot
			isScreenshotTakenChan <- true
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
