package rgb

import (
	"image"

	colorful "github.com/lucasb-eyer/go-colorful"
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

// Saturate returns a saturated copy of the given image by the given percent
// modify should be in range -1.0 to 1.0
func Saturate(img image.Image, modify float64) image.Image {
	if modify > 1.0 {
		modify = 1.0
	}
	if modify < -1.0 {
		modify = -1.0
	}
	sImg := image.NewRGBA(img.Bounds())
	size := img.Bounds().Size()
	for i := 0; i < size.X; i++ {
		for j := 0; j < size.Y; j++ {
			col, ok := colorful.MakeColor(img.At(i, j))
			if !ok {
				sImg.Set(i, j, col)
				continue
			}
			h, s, v := col.Hsv()
			var newSat, factor float64
			if modify >= 0 {
				room := 1 - s
				grey := s
				if s > 0.5 {
					grey = 1
				}
				factor = modify * room * grey
			} else {
				room := s
				factor = modify * room
			}
			newSat = s + factor
			sImg.Set(i, j, colorful.Hsv(h, newSat, v))
		}
	}
	return sImg
}
