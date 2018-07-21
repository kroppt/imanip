package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/kroppt/imanip/pkg/rgb"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Pass the filepath of image and saturation level")
		return
	}
	if len(os.Args) > 3 {
		fmt.Println("Too many arguments")
		return
	}
	inStr := os.Args[1]
	fIn, err := os.Open(inStr)
	defer fIn.Close()
	if err != nil {
		fmt.Println("File could not be opened")
		return
	}
	img, format, err := image.Decode(fIn)
	if err != nil {
		fmt.Println("Image type could not be decoded")
		return
	}
	s, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Println("Could not parse saturation level")
		return
	}
	gImg := rgb.Saturate(img, s)
	i := strings.LastIndexByte(inStr, '.')
	outStr := inStr[:i] + "_sat" + inStr[i:]
	fOut, err := os.OpenFile(outStr, os.O_CREATE|os.O_WRONLY, 0644)
	defer fOut.Close()
	if err != nil {
		fmt.Println("Could not open file", outStr)
		return
	}
	switch format {
	case "jpeg":
		// do not compress further
		err = jpeg.Encode(fOut, gImg, &jpeg.Options{Quality: 100})
		if err != nil {
			fmt.Println("File could not be encoded as jpeg")
		}
		return
	case "png":
		err = png.Encode(fOut, gImg)
		if err != nil {
			fmt.Println("File could not be encoded as png")
		}
		return
	default:
		fmt.Println("File format not recognized")
		return
	}
}
