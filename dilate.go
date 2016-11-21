package dilation

import (
	"image"
	"image/color"
	"image/draw"
)

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
	if alpha != 65535.0 {
		currentColor := dstImg.At(x, y)
		_, _, _, currentA := currentColor.RGBA()
		_, _, _, dstA := dstColor.RGBA()
		if currentA < dstA {
			targetColor := mixColors(srcColor, dstColor)
			dstImg.Set(x, y, targetColor)
		}
	}
}

// Should do something better here
func mixColors(c1, c2 color.Color) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	r := uint8((r1 + r2) / 65535.0 * 255)
	g := uint8((g1 + g2) / 65535.0 * 255)
	b := uint8((b1 + b2) / 65535.0 * 255)
	a := uint8((a1 + a2) / 65535.0 * 255)

	return color.RGBA{
		r, g, b, a,
	}
}
