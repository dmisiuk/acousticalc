package recording

import (
	"fmt"
	"image/png"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/imgo"
)

// TakeScreenshot takes a screenshot and saves it to the specified path.
func TakeScreenshot(path string) error {
	img := robotgo.CaptureScreen()
	return imgo.Save(path, img, png.DefaultCompression)
}

// GetScreenSize returns the screen size.
func GetScreenSize() (int, int) {
	return robotgo.GetScreenSize()
}

// MoveMouse moves the mouse to the specified coordinates.
func MoveMouse(x, y int) {
	robotgo.Move(x, y)
}

// ClickMouse clicks the mouse at the current position.
func ClickMouse() {
	robotgo.Click()
}

// ScrollMouse scrolls the mouse wheel.
func ScrollMouse(magnitude int, direction string) {
	robotgo.Scroll(magnitude, direction)
}

// TypeString types the given string.
func TypeString(text string) {
	robotgo.TypeStr(text)
}

// KeyTap taps the specified key.
func KeyTap(key string, modifiers ...interface{}) {
	robotgo.KeyTap(key, modifiers...)
}

// GetPixelColor returns the color of the pixel at the specified coordinates.
func GetPixelColor(x, y int) string {
	return robotgo.GetPixelColor(x, y)
}

// FindBitmap finds a bitmap on the screen.
func FindBitmap(path string) (int, int) {
	bitmap, _, err := robotgo.DecodeImg(path)
	if err != nil {
		return -1, -1
	}
	return robotgo.FindBitmap(bitmap)
}

// GetActiveWindowTitle returns the title of the active window.
func GetActiveWindowTitle() string {
	title, err := robotgo.GetTitle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting active window title: %v\n", err)
		return ""
	}
	return title
}
