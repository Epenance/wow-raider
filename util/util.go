package util

import (
	"fmt"
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
