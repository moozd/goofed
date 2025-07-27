package gfx

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Sdf struct {
	src  [][]byte
	buff [][]float32
	img  *image.Gray
}

func NewSdf(src [][]byte) *Sdf {
	s := &Sdf{
		src: src,
	}

	s.generate()

	return s
}

func (s *Sdf) Image() *image.Gray {
	return s.img
}

func (s *Sdf) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, s.img)
}

func (s *Sdf) generate() {
	height := len(s.src)
	width := len(s.src[0])

	buff := make([][]float32, height)
	for y := range height {
		buff[y] = make([]float32, width)
	}

	for y := range height {
		for x := range width {
			buff[y][x] = float32(s.sign(x, y) * s.computeNearestDistance(x, y))
		}
	}

	s.buff = buff

	s.img = s.createImage()
}

func (s *Sdf) sign(x, y int) float64 {
	if s.src[y][x] == 1 {
		return -1.0
	}

	return +1.0
}

func (s *Sdf) computeNearestDistance(x, y int) float64 {

	res := math.MaxFloat64

	coords := []struct{ stepX, stepY int }{
		//y x
		{0, 1},
		{0, -1},
		{1, 0},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, 0},
		{-1, -1},
	}

	vdc := 0

	for _, c := range coords {
		d, ok := s.findNearestOppositeNode(x, y, c.stepX, c.stepY)
		if ok && d < res {
			vdc += 1
			res = d
		}
	}

	if vdc == 0 {
		return 0
	}

	return res
}

func (s *Sdf) findNearestOppositeNode(x, y, stepX, stepY int) (d float64, ok bool) {
	maxX, maxY := len(s.src[0]), len(s.src)

	cx, cy := x+stepX, y+stepY
	for 0 <= cx && cx < maxX && 0 <= cy && cy < maxY {

		if s.src[cy][cx] != s.src[y][x] {
			d, ok = euclidean(x, y, cx, cy), true
			return

		}

		cx += stepX
		cy += stepY
	}

	d, ok = 0, false
	return
}

func (s *Sdf) createImage() *image.Gray {
	height := len(s.buff)
	width := len(s.buff[0])

	minVal, maxVal := float32(math.MaxFloat32), float32(-math.MaxFloat32)
	for y := range height {
		for x := range width {
			val := s.buff[y][x]
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
		}
	}

	img := image.NewGray(image.Rect(0, 0, width, height))
	scale := float32(255.0) / (maxVal - minVal)

	for y := range height {
		for x := range width {
			normalized := (s.buff[y][x] - minVal) * scale
			img.SetGray(x, y, color.Gray{Y: uint8(normalized)})
		}
	}

	return img
}

func euclidean(x0, y0, x1, y1 int) float64 {
	d := math.Sqrt(math.Pow(float64(x1-x0), 2) + math.Pow(float64(y1-y0), 2))
	return d
}
