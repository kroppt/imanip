package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/kroppt/imanip/pkg/rgb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Pass the filepath of file to gray as argument")
		return
	}
	if len(os.Args) > 2 {
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
	gImg := rgb.GrayImage(img)
	i := strings.LastIndexByte(inStr, '.')
	outStr := inStr[:i] + "_gray" + inStr[i:]
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
