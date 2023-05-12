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
		"{Screenshot: DisplayID: %v size: %vx%v}",
		s.DisplayID, s.Image.Rect.Max.X, s.Image.Rect.Max.Y)

}

func CompareScreenshots(s1, s2 *Screenshot) bool {
	// Check if the screenshots are for the same display
	if s1.DisplayID != s2.DisplayID {
		return false
	}

	// Check if the images have the same dimensions
	bounds1 := s1.Image.Bounds()
	bounds2 := s2.Image.Bounds()
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return false
	}

	// Check if the pixel values are the same
	// This is done by iterating over each pixel and comparing their RGBA values
	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			if s1.Image.RGBAAt(x, y) != s2.Image.RGBAAt(x, y) {
				return false
			}
		}
	}

	return true
}
