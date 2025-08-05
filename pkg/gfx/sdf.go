package gfx

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Pixel struct {
	x, y int
	v    byte
}

func (p *Pixel) GetKey() string {
	return fmt.Sprintf("%d|%d", p.x, p.y)
}

func (p *Pixel) GetDistance(t *Pixel) float64 {
	return math.Sqrt(math.Pow(float64(t.x-p.x), 2) + math.Pow(float64(t.y-p.y), 2))
}

type Bitmap struct {
	src  []byte
	h, w int
}

func (b *Bitmap) Get(x, y int) byte {

	if b.src[b.idx(x, y)] > 0 {
		return 1
	}
	return 0
}

func (b *Bitmap) idx(x, y int) int {
	return y*b.w + x
}

func (b *Bitmap) IsInRange(x, y int) bool {
	i := b.idx(x, y)

	return i > 0 && i < len(b.src)
}

func (b *Bitmap) Bounds() (height, width int) {
	height = b.h
	width = b.w
	return
}

type Sdf struct {
	input  Bitmap
	output [][]float32
	img    *image.Gray
}

func NewSdf(src []byte, h, w int) *Sdf {

	s := &Sdf{
		input: Bitmap{
			src: src,
			h:   h,
			w:   w,
		},
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
	height, width := s.input.Bounds()

	buff := make([][]float32, height)
	for y := range height {
		buff[y] = make([]float32, width)
	}

	for y := range height {
		for x := range width {
			buff[y][x] = float32(s.sign(x, y) * s.computeNearestOppositeDistance(x, y))
		}
	}

	fmt.Print("[\n")
	for _, row := range buff {
		fmt.Printf("\t%v\n", row)
	}
	fmt.Println("]")
	fmt.Println()

	s.output = buff

	s.img = s.createImage()
}

func (s *Sdf) sign(x, y int) float64 {
	if s.input.Get(x, y) == 1 {
		return -1.0
	}

	return +1.0
}

func (s *Sdf) computeNearestOppositeDistance(x, y int) float64 {
	var node *Pixel
	found := false
	root := &Pixel{x, y, s.input.Get(x, y)}
	visited := make(map[string]bool)
	visited[root.GetKey()] = true

	queue := []*Pixel{root}
	vx := []int{0, 0, 1, 1, 1, -1, -1, -1}
	vy := []int{1, -1, 0, 1, -1, 1, 0, -1}

	for len(queue) > 0 {
		node = queue[0]
		queue = queue[1:]

		if node != root && node.v != root.v {
			found = true
			break
		}

		for i := range len(vy) {
			dx := node.x + vx[i]
			dy := node.y + vy[i]

			if !s.input.IsInRange(dx, dy) {
				continue
			}

			p := &Pixel{dx, dy, s.input.Get(dx, dy)}
			if _, ok := visited[p.GetKey()]; ok {
				continue
			}
			queue = append(queue, p)
			visited[p.GetKey()] = true

		}

	}

	if !found {
		return 0
	}

	return node.GetDistance(root)

}

func (s *Sdf) createImage() *image.Gray {
	height := len(s.output)
	width := len(s.output[0])

	minVal, maxVal := float32(math.MaxFloat32), float32(-math.MaxFloat32)
	for y := range height {
		for x := range width {
			val := s.output[y][x]
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
			normalized := (s.output[y][x] - minVal) * scale
			img.SetGray(x, y, color.Gray{Y: uint8(normalized)})
		}
	}

	return img
}
