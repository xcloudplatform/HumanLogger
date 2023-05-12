package main

import (
	"fmt"
	"github.com/ClickerAI/ClickerAI/adapters/logger"
	"github.com/ClickerAI/ClickerAI/adapters/packer"
	"github.com/ClickerAI/ClickerAI/adapters/uploader"
	"github.com/ClickerAI/ClickerAI/core"
	"os"
)

func main() {
	c := core.NewCore()

	l := logger.NewLocalFileLogger()
	p := packer.ZipPacker{}
	u := uploader.GoogleStorageUploader{BucketName: "sessions-upload-v1"}
	session, _ := l.StartLogging()

	go func() {
		for ev := range c.UIEventStream {
			err := l.LogUIEvent(session, &ev)
			if err != nil {
				fmt.Printf("error logging event: %v", err)
			}
		}
	}()

	go func() {
		for scr := range c.ScreenshotStream {
			err := l.LogScreenshot(session, &scr)
			if err != nil {
				fmt.Printf("error writing screenshot: %v", err)
			}
		}
	}()

	go func() {
		for active := range c.UserActivity.IsActiveStream {
			if !active && session.IsNeedsRotation() {
				err := l.StopLogging(session)
				if err != nil {
					fmt.Printf("error stoping logging: %v", err)

				}
				go func() {
					packedSession, err := p.Pack(session)
					if err == nil {
						fmt.Printf("packed session: %v\n", packedSession)
						err := u.Upload(&session, packedSession)
						if err != nil {
							fmt.Printf("error uploading: %v\n", err)
						} else {
							err := os.Remove(packedSession)
							if err != nil {
								fmt.Printf("Error deleting file: %v\n", err)
							}

						}
					}
				}()

				c.ResetLastScreenshots()

				session, err = l.StartLogging()
				if err != nil {
					fmt.Printf("error starting logging: %v", err)

				}

			}
		}
	}()

	err := c.Start()
	if err != nil {
		return
	}

}
