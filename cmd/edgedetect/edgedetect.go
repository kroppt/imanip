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
		fmt.Println("Pass the filepath of image to process and the number of " +
			"the operator, and optionally the threshold for cutoff")
		fmt.Println("1 for Sobel-Feldman")
		fmt.Println("2 for Scharr")
		fmt.Println("3 for Prewitt")
		return
	}
	if len(os.Args) > 4 {
		fmt.Println("Too many arguments")
		return
	}
	inStr := os.Args[1]
	fIn, err := os.Open(inStr)
	defer fIn.Close()
	if err != nil {
		fmt.Println("File could not be opened:", err)
		return
	}
	img, format, err := image.Decode(fIn)
	if err != nil {
		fmt.Println("Image type could not be decoded:", err)
		return
	}
	opInd, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Could not parse operator index as integer")
		return
	}
	thresh := -1.0
	if len(os.Args) > 3 {
		thresh, err = strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Println("Could not parse threshold as float64")
			return
		}
	}
	var op rgb.EdgeOperator
	switch opInd {
	case 1:
		op = rgb.NewSobelFeldman()
		break
	case 2:
		op = rgb.NewScharr()
		break
	case 3:
		op = rgb.NewPrewitt()
		break
	default:
		fmt.Println("")
		return
	}
	gImg, err := rgb.EdgeDetect(img, op, thresh)
	if err != nil {
		fmt.Println("Could not process:", err)
		return
	}
	i := strings.LastIndexByte(inStr, '.')
	outStr := inStr[:i] + "_edges" + inStr[i:]
	fOut, err := os.OpenFile(outStr, os.O_CREATE|os.O_WRONLY, 0644)
	defer fOut.Close()
	if err != nil {
		fmt.Println("Could not open file:", err)
		return
	}
	switch format {
	case "jpeg":
		// do not compress further
		err = jpeg.Encode(fOut, gImg, &jpeg.Options{Quality: 100})
		if err != nil {
			fmt.Println("File could not be encoded as jpeg:", err)
		}
		return
	case "png":
		err = png.Encode(fOut, gImg)
		if err != nil {
			fmt.Println("File could not be encoded as png:", err)
		}
		return
	default:
		fmt.Println("File format not recognized")
		return
	}
}
