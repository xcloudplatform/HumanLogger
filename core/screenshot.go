package core

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"image"
	"time"
)

type Screenshot struct {
	Timestamp time.Time

	Image     *image.RGBA
	DisplayID int
}

func CaptureScreenshots() ([]Screenshot, error) {
	num := screenshot.NumActiveDisplays()

	screenshots := make([]Screenshot, num)

	for i := 0; i < num; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			continue
		}
		screenshot := Screenshot{
			Timestamp: time.Now(),
			Image:     img,
			DisplayID: i,
		}

		screenshots[i] = screenshot
	}

	return screenshots, nil
}

func (s Screenshot) String() string {
	return fmt.Sprintf(
		"{Screenshot: DisplayID: %v}",
		s.DisplayID)

}
