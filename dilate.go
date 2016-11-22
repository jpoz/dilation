package dilation

import (
	"image"
	"image/color"
	"image/draw"
)

// MAXC is the maxium value returned by RGBA() will have
const MAXC = 65535.0

// EditableImage is an Image but allows color to be set
// image.Alpha, image.Alpha16, image.CMYK, image.Gray
// image.Gray16, image.NRGBA, image.Paletted, etc. all conform to this
// interface.
type EditableImage interface {
	Set(x, y int, c color.Color)
	image.Image
}

// Dialate will dialate a given image around transparent edges
func Dialate(dstImg EditableImage) error {
	b := dstImg.Bounds()
	origImg := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(origImg, origImg.Bounds(), dstImg, b.Min, draw.Src)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			err := dialatePoint(origImg, dstImg, x, y)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func dialatePoint(origImg image.Image, dstImg EditableImage, x, y int) error {
	srcColor := origImg.At(x, y)
	_, _, _, alpha := srcColor.RGBA()
	if alpha != 0 {
		// TODO setup config
		for r := 1; r < 10; r++ {
			drawCircle(origImg, dstImg, x, y, r, uint8(alpha))
		}
	}
	return nil
}

func drawCircle(origImg image.Image, dstImg EditableImage, x, y, r int, a uint8) {
	x1, y1, rad := -r, 0, 2-2*r
	for {
		drawDialatePixel(origImg, dstImg, x-x1, y+y1, a)
		drawDialatePixel(origImg, dstImg, x-y1, y-x1, a)
		drawDialatePixel(origImg, dstImg, x+x1, y-y1, a)
		drawDialatePixel(origImg, dstImg, x+y1, y+x1, a)
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

func drawDialatePixel(origImg image.Image, dstImg EditableImage, x, y int, a uint8) {
	r, g, b, _ := color.Black.RGBA()
	srcColor := origImg.At(x, y)
	dstColor := color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: a,
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

	r := uint8(mixColor(r1, r2, a1, a2) / MAXC * 255)
	g := uint8(mixColor(g1, g2, a1, a2) / MAXC * 255)
	b := uint8(mixColor(b1, b2, a1, a2) / MAXC * 255)
	// This should not be needed
	a := uint8(mixColor(a1, a2, a1, a2) / MAXC * 255)

	return color.RGBA{
		r, g, b, a,
	}
}

func mixColor(cv1, cv2, av1, av2 uint32) float64 {
	if av1 == 0 && av2 == 0 {
		return 0.0
	}
	a1 := float64(av1) / MAXC
	a2 := float64(av2) / MAXC
	c1 := float64(cv1)
	c2 := float64(cv2)

	a0 := (a1 + a2*(1-a1))
	c0 := (c1*a1 + c2*a2*(1-a1)) / a0

	return c0
}
