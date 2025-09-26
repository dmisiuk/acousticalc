//go:build visualtests

package visual

import (
	"fmt"
	"image"
	"runtime"

	"github.com/go-vgo/robotgo"
)

type RobotGoEngine struct {
	platform string
}

func NewRobotGoEngine() ScreenshotEngine {
	return &RobotGoEngine{platform: runtime.GOOS}
}

func (rg *RobotGoEngine) Capture() (image.Image, error) {
	bitmap := robotgo.CaptureScreen()
	if bitmap == nil {
		return nil, fmt.Errorf("failed to capture screen with robotgo")
	}
	defer robotgo.FreeBitmap(bitmap)

	img := robotgo.ToImage(bitmap)
	if img == nil {
		return nil, fmt.Errorf("failed to convert bitmap to image")
	}
	return img, nil
}

func (rg *RobotGoEngine) GetPlatform() string {
	return rg.platform
}

func (rg *RobotGoEngine) IsAvailable() bool {
	return true
}
