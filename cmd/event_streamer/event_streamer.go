package main

import (
	"github.com/ClickerAI/ClickerAI/adapters/logger"
	"github.com/ClickerAI/ClickerAI/adapters/packer"
	"github.com/ClickerAI/ClickerAI/adapters/uploader"
	"github.com/ClickerAI/ClickerAI/core"
)

func main() {
	c := core.NewCore()

	l := logger.NewLocalFileLogger()
	p := packer.ZipPacker{}
	u := uploader.S3Uploader{}
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

	go func() {
		for active := range c.UserActivity.IsActiveStream {
			if !active && session.IsNeedsRotation() {
				l.StopLogging(session)
				packedSession, err := p.Pack(session)
				if err == nil {
					u.Upload(&session, packedSession)
				}
				c.ResetLastScreenshots()

				session, _ = l.StartLogging()

			}
		}
	}()

	c.Start()

}
