package rgb

import (
	"image"
)

// GrayImage returns a pixel-gray copy of the given image
func GrayImage(img image.Image) image.Image {
	gImg := image.NewGray(img.Bounds())
	size := img.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			g := gImg.ColorModel().Convert(img.At(x, y))
			gImg.Set(x, y, g)
		}
	}
	return gImg
}
