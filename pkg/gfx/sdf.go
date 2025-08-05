package gfx

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type sdf struct {
	src bitmap
}

func generateSDF(glyph []byte, h, w int) *image.Gray {

	s := &sdf{
		src: bitmap{
			bytes: glyph,
			h:     h,
			w:     w,
		},
	}

	height, width := s.src.bounds()

	buff := make([][]float32, height)
	for y := range height {
		buff[y] = make([]float32, width)
	}

	for y := range height {
		for x := range width {
			buff[y][x] = s.sign(x, y) * s.computeNearestOppositeDistance(x, y)
		}
	}

	return s.createImage(buff)
}

func (s *sdf) sign(x, y int) float32 {
	if s.src.get(x, y) == 1 {
		return -1.0
	}

	return +1.0
}

func (s *sdf) computeNearestOppositeDistance(x, y int) float32 {
	// this the current node
	var node *pixel
	found := false
	root := &pixel{x, y, s.src.get(x, y)}
	visited := make(map[string]bool)
	visited[root.getKey()] = true

	queue := []*pixel{root}
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

			if !s.src.isInRange(dx, dy) {
				continue
			}

			p := &pixel{dx, dy, s.src.get(dx, dy)}

			if _, ok := visited[p.getKey()]; ok {
				continue
			}

			queue = append(queue, p)
			visited[p.getKey()] = true

		}

	}

	if !found {
		return 0
	}

	return node.measure(root)

}

func (s *sdf) createImage(output [][]float32) *image.Gray {
	height := len(output)
	width := len(output[0])

	minVal, maxVal := float32(math.MaxFloat32), float32(-math.MaxFloat32)
	for y := range height {
		for x := range width {
			val := output[y][x]
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
			normalized := (output[y][x] - minVal) * scale
			img.SetGray(x, y, color.Gray{Y: uint8(normalized)})
		}
	}

	return img
}

type pixel struct {
	x, y int
	v    byte
}

func (p *pixel) getKey() string {
	return fmt.Sprintf("%d|%d", p.x, p.y)
}

func (p *pixel) measure(t *pixel) float32 {
	return float32(math.Sqrt(math.Pow(float64(t.x-p.x), 2) + math.Pow(float64(t.y-p.y), 2)))
}

type bitmap struct {
	bytes []byte
	h, w  int
}

func (b *bitmap) get(x, y int) byte {

	if b.bytes[b.idx(x, y)] > 0 {
		return 1
	}
	return 0
}

func (b *bitmap) idx(x, y int) int {
	return y*b.w + x
}

func (b *bitmap) isInRange(x, y int) bool {
	i := b.idx(x, y)

	return i > 0 && i < len(b.bytes)
}

func (b *bitmap) bounds() (height, width int) {
	height = b.h
	width = b.w
	return
}
