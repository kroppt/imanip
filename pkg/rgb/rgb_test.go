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

func TestSaturate(t *testing.T) {
	iImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	col := color.RGBA{0x00, 0x00, 0xff, 0xff}
	iImg.SetRGBA(0, 0, col)
	wImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	col = color.RGBA{0xff, 0xff, 0xff, 0xff}
	wImg.SetRGBA(0, 0, col)
	zImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	col = color.RGBA{0xff, 0xff, 0xff, 0x00}
	zImg.SetRGBA(0, 0, col)
	bImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	col = color.RGBA{0x00, 0x00, 0x00, 0xff}
	bImg.SetRGBA(0, 0, col)
	type args struct {
		img    image.Image
		modify float64
	}
	tests := []struct {
		name string
		args args
		want image.Image
	}{
		{"OverflowSaturate", args{iImg, 5.0}, iImg},
		{"IdentitySaturate", args{iImg, 1.0}, iImg},
		{"IdentityDesaturate", args{iImg, -1.0}, wImg},
		{"UnderflowDesaturate", args{iImg, -5.0}, wImg},
		{"ZeroAlpha", args{zImg, 0.5}, bImg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Saturate(tt.args.img, tt.args.modify); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Saturate() = %v, want %v", got, tt.want)
			}
		})
	}
}
