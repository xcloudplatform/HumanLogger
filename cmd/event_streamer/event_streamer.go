package main

import (
	"github.com/ClickerAI/ClickerAI/adapters/logger"
	"github.com/ClickerAI/ClickerAI/core"
)

func main() {
	c := core.NewCore()

	l := logger.NewLocalFileLogger()

	session, _ := l.StartLogging()

	go func() {
		for ev := range c.UIEventStream {
			l.LogUIEvent(session, &ev)
		}
	}()

	go func() {
		for scr := range c.ScreenshotStream {
			l.LogScreenshot(session, &scr)
		}
	}()

	c.Start()

}
