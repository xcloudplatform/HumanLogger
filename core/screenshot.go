package core

import (
	"fmt"
	"github.com/kbinani/screenshot"

	"image"
	"image/color"
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
		scr := Screenshot{
			Timestamp: time.Now(),
			Image:     img,
			DisplayID: i,
		}

		fmt.Printf("%d captured screenshot\n", scr.Timestamp.UnixNano())
		screenshots[i] = scr
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

func (s Screenshot) Diff(other *Screenshot) (bool, *Screenshot, error) {
	// check if screenshots have the same dimensions
	if !s.Image.Bounds().Eq(other.Image.Bounds()) {
		return false, nil, fmt.Errorf("screenshots have different dimensions")
	}

	// create a new RGBA image with the same dimensions
	bounds := s.Image.Bounds()
	diffImg := image.NewRGBA(bounds)

	// flag to track whether the screenshots are the same
	same := true

	// loop through each pixel and compare them
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := s.Image.At(x, y)
			c2 := other.Image.At(x, y)

			// if the pixels are the same, set it to magenta with 0 alpha
			if c1 == c2 {
				diffImg.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0})
			} else {
				diffImg.Set(x, y, s.Image.At(x, y))
				same = false
			}
		}
	}

	// check if the screenshots are the same
	if same {
		return true, nil, nil
	}

	// create a new screenshot with metadata from the new screenshot and the diff image
	diffScreenshot := &Screenshot{
		Timestamp: s.Timestamp,
		Image:     diffImg,
		DisplayID: s.DisplayID,
	}

	return false, diffScreenshot, nil
}
