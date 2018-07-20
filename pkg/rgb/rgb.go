package rgb

import (
	"image"
)

// GrayImage turns any image into its pixel-gray counterpart
func GrayImage(img image.Image) image.Image {
	gray := &image.Gray{}
	gImg := image.NewGray(img.Bounds())
	size := img.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			g := gray.ColorModel().Convert(img.At(x, y))
			gImg.Set(x, y, g)
		}
	}
	return gImg
}
