package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/jpoz/dilation"
)

func main() {
	f, err := os.Open("./examples/big.png")
	check(err)

	src, err := png.Decode(f)
	check(err)

	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)

	config := dilation.DialateConfig{
		Stroke:      10,
		StrokeColor: color.Black,
	}
	dilation.Dialate(img, config)

	f2, err := os.Create("./big-output.png")
	check(err)

	err = png.Encode(f2, img)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
