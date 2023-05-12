package core

import (
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
