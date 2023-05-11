package main

import (
	"fmt"
	"time"
	// "time"

	hook "github.com/robotn/gohook"
)

func main() {
	isActiveChan := make(chan bool, 200)
	go handleEvents(isActiveChan)
	go takeScreenshots(isActiveChan)
	go getWindowTitles()
	select {}
}


func takeScreenshots(isActiveChan chan bool) {
	isActive := false
	for {
		select {
		case value, ok := <-isActiveChan:
			if ok {
				isActive = value
			} else {
				isActive = false
			}
		default:
			isActive = false
		}

		interval := 5 * time.Second
		if isActive {
			interval = 300 * time.Millisecond
		}
		fmt.Println("Sleeping for", interval)
		time.Sleep(interval)
		
		fmt.Println("TAKING SCREENSHOT!")
	}
}


func getWindowTitles() {

}

func handleEvents(isActiveChan chan bool) {
	lastEventTime := time.Now()
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
		isActiveChan <- true
		lastEventTime = time.Now()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if time.Since(lastEventTime) > 3*time.Second {
				isActiveChan <- false
			}
		}
	}
}
