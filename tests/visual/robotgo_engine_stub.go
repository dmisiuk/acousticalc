//go:build !visualtests

package visual

import (
	"errors"
	"image"
)

func NewRobotGoEngine() ScreenshotEngine {
	return &RobotGoEngine{}
}

type RobotGoEngine struct{}

func (rg *RobotGoEngine) Capture() (image.Image, error) {
	return nil, errors.New("robotgo engine unavailable; rebuild with -tags visualtests")
}

func (rg *RobotGoEngine) GetPlatform() string {
	return "unavailable"
}

func (rg *RobotGoEngine) IsAvailable() bool {
	return false
}
