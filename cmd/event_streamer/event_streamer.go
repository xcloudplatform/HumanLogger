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
			err := l.LogUIEvent(session, &ev)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		for scr := range c.ScreenshotStream {
			err := l.LogScreenshot(session, &scr)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		for active := range c.UserActivity.IsActiveStream {
			if !active && session.IsNeedsRotation() {
				err := l.StopLogging(session)
				if err != nil {
					return
				}
				packedSession, err := p.Pack(session)
				if err == nil {
					err := u.Upload(&session, packedSession)
					if err != nil {
						return
					}
				}
				c.ResetLastScreenshots()

				session, _ = l.StartLogging()

			}
		}
	}()

	err := c.Start()
	if err != nil {
		return
	}

}
