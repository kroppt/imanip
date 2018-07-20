package rgb_test

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	. "github.com/kroppt/imanip/pkg/rgb"
)

func TestGrayImage(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 10, 10))
	gImg := image.NewGray(image.Rect(0, 0, 10, 10))
	cImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	cgImg := image.NewGray(image.Rect(0, 0, 10, 10))
	{
		col := color.RGBA{0x00, 0x00, 0xff, 0xff}
		gcol, _, _, _ := gImg.ColorModel().Convert(col).RGBA()
		cImg.SetRGBA(0, 0, col)
		cgImg.SetGray(0, 0, color.Gray{uint8(gcol)})
	}
	tests := []struct {
		name string
		arg  image.Image
		want image.Image
	}{
		{"GrayToGray", img, gImg},
		{"ColorToGray", cImg, cgImg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GrayImage(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrayImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
