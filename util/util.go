package util

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"image"
	"io"
	"os"
	"time"
)

var writer io.Writer = os.Stdout

func PixelColorAt(img image.Image, x int, y int) (r, g, b uint32) {
	c := img.At(x, y)
	r, g, b, _ = c.RGBA()
	// Convert the color's components from 16-bit to 8-bit
	return r / 257, g / 257, b / 257
}

func SetWriter(w io.Writer) {
	writer = w
}

func Log(msg string) {
	t := time.Now()

	fmt.Fprintln(writer, "["+t.Format("15:04:05")+"] "+msg)
}

type RGB struct {
	R uint32
	G uint32
	B uint32
}

var RED = RGB{R: 255, G: 2, B: 0}
var BLUE = RGB{R: 0, G: 65, B: 255}
var PURPLE = RGB{R: 150, G: 0, B: 255}
var GREEN = RGB{R: 0, G: 255, B: 3}
var YELLOW = RGB{R: 255, G: 225, B: 0}

// Color map called COLORS for exports of all our colors
var COLORS = map[string]RGB{
	"RED":    RED,
	"BLUE":   BLUE,
	"PURPLE": PURPLE,
	"GREEN":  GREEN,
	"YELLOW": YELLOW,
}

func IsColor(rgb RGB, img image.Image, x int, y int, debug ...bool) bool {
	var isDebug bool
	if len(debug) > 0 {
		isDebug = debug[0]
	} else {
		isDebug = false // Default value
	}

	r, g, b := PixelColorAt(img, x, y)

	correctColors := r == rgb.R && g == rgb.G && b == rgb.B

	if isDebug {
		fmt.Printf("Matches: %d R: %d, G: %d, B: %d\n", correctColors, r, g, b)
	}

	return correctColors
}

func GetColor(b bool, c1 tcell.Color, c2 tcell.Color) tcell.Color {
	if b {
		return c1
	}
	return c2
}
