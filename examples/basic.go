package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/jpoz/dilation"
)

func main() {
	f, err := os.Open("./example.png")
	check(err)

	src, err := png.Decode(f)
	check(err)

	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)

	dilation.Dialate(img)

	f2, err := os.Create("./example-output.png")
	check(err)

	err = png.Encode(f2, img)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
