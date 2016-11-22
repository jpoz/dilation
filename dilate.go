package dilation

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// MAXC is the maxium value returned by RGBA() will have
const MAXC = 65535.0

// DialateConfig settings for the dialation
type DialateConfig struct {
	Stroke      int
	StrokeColor color.Color
	// TODO make feather option
	// Feather     int
}

// Dialate will dialate a given image around transparent edges
func Dialate(dstImg draw.Image, config DialateConfig) error {
	b := dstImg.Bounds()
	origImg := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(origImg, origImg.Bounds(), dstImg, b.Min, draw.Src)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			err := dialatePoint(origImg, dstImg, x, y, config)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func dialatePoint(origImg image.Image, dstImg draw.Image, x, y int, config DialateConfig) error {
	srcColor := origImg.At(x, y)
	_, _, _, alpha := srcColor.RGBA()
	if alpha != 0 {
		// TODO setup config
		for r := 1; r <= config.Stroke; r++ {
			drawCircle(origImg, dstImg, x, y, r, alpha, config)
		}
	}
	return nil
}

func drawCircle(origImg image.Image, dstImg draw.Image, x, y, r int, a uint32, config DialateConfig) {
	x1, y1, rad := -r, 0, 2-2*r
	for {
		drawDialatePixel(origImg, dstImg, x-x1, y+y1, a, config)
		drawDialatePixel(origImg, dstImg, x-y1, y-x1, a, config)
		drawDialatePixel(origImg, dstImg, x+x1, y-y1, a, config)
		drawDialatePixel(origImg, dstImg, x+y1, y+x1, a, config)
		r = rad
		if r > x1 {
			x1++
			rad += x1*2 + 1
		}
		if r <= y1 {
			y1++
			rad += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

func drawDialatePixel(origImg image.Image, dstImg draw.Image, x, y int, a uint32, config DialateConfig) {
	r, g, b, _ := config.StrokeColor.RGBA()
	srcColor := origImg.At(x, y)
	dstColor := color.RGBA{
		R: cnv(r),
		G: cnv(g),
		B: cnv(b),
		A: cnv(a),
	}
	_, _, _, alpha := srcColor.RGBA()
	if alpha != MAXC {
		currentColor := dstImg.At(x, y)
		_, _, _, currentA := currentColor.RGBA()
		_, _, _, dstA := dstColor.RGBA()
		if currentA < dstA {
			targetColor := mixColors(srcColor, dstColor)
			dstImg.Set(x, y, targetColor)
		}
	}
}

// Using Alpha compositing:  https://en.wikipedia.org/wiki/Alpha_compositing
func mixColors(c1, c2 color.Color) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	r := mixColor(r1, r2, a1, a2)
	g := mixColor(g1, g2, a1, a2)
	b := mixColor(b1, b2, a1, a2)
	// This should not be needed, the alpha value is computed in each ofthe
	// mixColor calls. TODO refactor
	a := mixColor(a1, a2, a1, a2)

	return color.RGBA{r, g, b, a}
}

func mixColor(cv1, cv2, av1, av2 uint32) uint8 {
	if av1 == 0 && av2 == 0 {
		return 0.0
	}
	a1 := float64(av1) / MAXC
	a2 := float64(av2) / MAXC
	c1 := float64(cv1)
	c2 := float64(cv2)

	a0 := (a1 + a2*(1-a1))
	c0 := (c1*a1 + c2*a2*(1-a1)) / a0
	// THis max might not be needed
	c0 = math.Max(0.0, math.Min(c0, MAXC))

	return uint8(c0 / MAXC * 255)
}

// util function to go between Golang stupid values
func cnv(c uint32) uint8 {
	return uint8(c / MAXC * 255)
}
