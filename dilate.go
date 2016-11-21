package dilation

import (
	"fmt"
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
	r := 1
	srcColor := origImg.At(x, y)
	_, _, _, alpha := srcColor.RGBA()
	fmt.Println(alpha)
	if alpha != 0 {
		x1, y1, rad := -r, 0, 2-2*r
		for {
			drawDialatePixel(origImg, dstImg, x-x1, y+y1)
			drawDialatePixel(origImg, dstImg, x-y1, y-x1)
			drawDialatePixel(origImg, dstImg, x+x1, y-y1)
			drawDialatePixel(origImg, dstImg, x+y1, y+x1)
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

	return nil
}

func drawDialatePixel(origImg image.Image, dstImg EditableImage, x, y int) {
	dialateColor := color.Black
	srcColor := origImg.At(x, y)
	_, _, _, alpha := srcColor.RGBA()
	if alpha == 0 {
		dstImg.Set(x, y, dialateColor)
	}
}
