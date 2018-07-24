package rgb

import (
	"errors"
	"image"
	"image/color"
	"math"

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

// EdgeOperator contains an averaging and a differentiation kernel
// The kernels are used in convolution for edge detection in images
// The threshold is the point at which low-value points are equalized
type EdgeOperator struct {
	AvgKern   []int
	DiffKern  []int
	Threshold float64
}

// NewSobelFeldman returns a Sobel-Feldman edge-detection operator
func NewSobelFeldman() EdgeOperator {
	return EdgeOperator{
		[]int{1, 2, 1},
		[]int{1, 0, -1},
		70.0,
	}
}

// NewScharr returns a Scharr edge-detection operator
func NewScharr() EdgeOperator {
	return EdgeOperator{
		[]int{3, 10, 3},
		[]int{1, 0, -1},
		70.0,
	}
}

// NewPrewitt returns a Prewitt edge-detection operator
func NewPrewitt() EdgeOperator {
	return EdgeOperator{
		[]int{1, 1, 1},
		[]int{1, 0, -1},
		70.0,
	}
}

// GradientX returns the full convolution matrix for the x-axis
func (op EdgeOperator) GradientX() ([][]int, error) {
	if len(op.AvgKern) == 0 || len(op.DiffKern) == 0 {
		return nil, errors.New("Edge Operator has kernels with zero length")
	}
	grad := make([][]int, len(op.AvgKern))
	for i := 0; i < len(grad); i++ {
		grad[i] = make([]int, len(op.DiffKern))
	}
	for y := 0; y < len(grad); y++ {
		for x := 0; x < len(grad[y]); x++ {
			grad[y][x] = op.AvgKern[y] * op.DiffKern[x]
		}
	}
	return grad, nil
}

// GradientY returns the full convolution matrix for the y-axis
func (op EdgeOperator) GradientY() ([][]int, error) {
	if len(op.AvgKern) == 0 || len(op.DiffKern) == 0 {
		return nil, errors.New("Edge Operator has kernels with zero length")
	}
	grad := make([][]int, len(op.DiffKern))
	for i := 0; i < len(grad); i++ {
		grad[i] = make([]int, len(op.AvgKern))
	}
	for y := 0; y < len(grad); y++ {
		for x := 0; x < len(grad[y]); x++ {
			grad[y][x] = op.AvgKern[x] * op.DiffKern[y]
		}
	}
	return grad, nil
}

// EdgeDetect returns an image with edges highlighted using the given operator
func EdgeDetect(img image.Image, op EdgeOperator, threshold float64) (image.Image, error) {
	eImg := image.NewGray(img.Bounds())
	if threshold < 0 {
		threshold = op.Threshold
	}
	if threshold < 0 || threshold > 0xff {
		return eImg, errors.New("Operator threshold out of 0-255 bounds")
	}
	xGrad, err := op.GradientX()
	if err != nil {
		return eImg, err
	}
	yGrad, err := op.GradientY()
	if err != nil {
		return eImg, err
	}
	if len(xGrad)%2 == 0 || len(xGrad)%2 == 0 {
		return eImg, errors.New("Gradient of non-odd length")
	}
	if len(yGrad)%2 == 0 || len(yGrad)%2 == 0 {
		return eImg, errors.New("Gradient of non-odd length")
	}
	if len(xGrad) != len(yGrad) || len(xGrad[0]) != len(yGrad[0]) {
		return eImg, errors.New("Gradients do not have equal size")
	}
	size := img.Bounds().Size()
	minX := len(xGrad[0])
	if len(yGrad[0]) < minX {
		minX = len(yGrad[0])
	}
	if size.X < minX {
		return eImg,
			errors.New("Image size is smaller than the gradient size")
	}
	minY := len(xGrad)
	if len(yGrad) < minX {
		minY = len(yGrad)
	}
	if size.Y < minY {
		return eImg,
			errors.New("Image size is smaller than the gradient size")
	}
	// can't convolute past edge of image
	// black sides of len(mat)/2 length
	xbuf := len(xGrad) / 2
	ybuf := len(xGrad[0]) / 2
	for y := ybuf; y < size.Y-ybuf; y++ {
		for x := xbuf; x < size.X-xbuf; x++ {
			var xSum float64
			for my := range xGrad {
				for mx, m := range xGrad[my] {
					c := img.At(x+xbuf-mx, y+ybuf-my)
					g, _, _, _ := eImg.ColorModel().Convert(c).RGBA()
					xSum += float64(int32(g>>8) * int32(m))
				}
			}
			var ySum float64
			for my := range yGrad {
				for mx, m := range yGrad[my] {
					c := img.At(x+xbuf-mx, y+ybuf-my)
					g, _, _, _ := eImg.ColorModel().Convert(c).RGBA()
					ySum += float64(int32(g>>8) * int32(m))
				}
			}
			mag := math.Sqrt(ySum*ySum + xSum*xSum)
			if mag < threshold {
				mag = threshold
			}
			// if magnitude exceeds 255
			if mag > float64(^uint8(0)) {
				mag = float64(^uint8(0))
			}
			eImg.SetGray(x, y, color.Gray{uint8(mag)})
		}
	}
	return eImg, nil
}

// Convolute returns a convoluted copy of the given image with the given matrix
func Convolute(img image.Image, mat [][]float64) (image.Image, error) {
	cImg := image.NewGray(img.Bounds())
	if len(mat) == 0 || len(mat[0]) == 0 {
		return cImg, errors.New("convolution matrix of zero length")
	}
	if len(mat)%2 == 0 || len(mat[0])%2 == 0 {
		return cImg, errors.New("convolution matrix of non-odd length")
	}
	size := img.Bounds().Size()
	min := size.X
	if size.Y < size.X {
		min = size.Y
	}
	if min < len(mat) {
		return cImg, errors.New("convolution matrix is smaller than " +
			"the size of the image")
	}
	var norm float64
	for my := 0; my < len(mat); my++ {
		for mx := 0; mx < len(mat[my]); mx++ {
			norm += mat[my][mx]
		}
	}
	if norm == 0 {
		norm = 1
	} else {
		norm = 1 / norm
	}
	// can't convolute past edge of image
	// black sides of len(mat)/2 length
	xbuf := len(mat) / 2
	ybuf := len(mat[0]) / 2
	for y := ybuf; y < size.Y-ybuf; y++ {
		for x := xbuf; x < size.X-xbuf; x++ {
			var sum float64
			for my := range mat {
				for mx, m := range mat[my] {
					c := img.At(x+xbuf-mx, y+ybuf-my)
					g, _, _, _ := cImg.ColorModel().Convert(c).RGBA()
					sum += float64(g>>8) * m
				}
			}
			sum *= norm
			if sum < 0 {
				sum *= -1
			}
			if sum > float64(^uint8(0)) {
				sum = float64(^uint8(0))
			}
			cImg.SetGray(x, y, color.Gray{uint8(sum)})
		}
	}
	return cImg, nil
}
