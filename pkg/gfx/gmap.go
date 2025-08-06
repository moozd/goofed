package gfx

import (
	"image"
	"image/draw"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type GMap struct {
	font         *opentype.Font
	face         font.Face
	cache        map[rune]*GData
	LineHeight   int
	AdvanceWidth int
}

type GData struct {
	char          rune
	Source        *image.Gray
	DistanceField *image.Gray
}

func NewGMap(addr string, size int) (*GMap, error) {

	gm := &GMap{
		cache: make(map[rune]*GData),
	}

	gm.createFace(addr, float32(size))
	gm.computeBoundingBoxSize()

	return gm, nil
}

func (gm *GMap) Close() {
	gm.face.Close()
}

func (gm *GMap) Get(r rune) (*GData, bool) {
	meta, ok := gm.cache[r]

	if ok {
		return meta, true
	}

	meta, ok = gm.create(r)
	if !ok {
		return nil, false
	}

	gm.cache[r] = meta

	return meta, true
}

func (gm *GMap) createFace(addr string, size float32) error {

	fb, err := os.ReadFile(addr)
	if err != nil {
		return err
	}

	fnt, err := opentype.Parse(fb)
	if err != nil {
		return err
	}

	fc, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     92,
		Hinting: font.HintingNone,
	})

	if err != nil {
		return err
	}

	gm.face = fc

	return nil
}

func (gm *GMap) computeBoundingBoxSize() {

	chars := []rune{'M', 'W', 'Q', '#', '%', '&'}
	width, height := 0, 0

	for _, char := range chars {

		_, adv, ok := gm.face.GlyphBounds(char)
		metrics := gm.face.Metrics()

		w := int(adv >> 6)            // font width
		h := int(metrics.Height >> 6) // font height

		if ok && w > width {
			width = w
		}

		if h > height {
			height = h
		}
	}

	log.Printf("LineHeight: %d , AdvanceWidth: %d", height, width)

	gm.LineHeight = height
	gm.AdvanceWidth = width
}

func (gm *GMap) create(r rune) (*GData, bool) {

	bounds, advance, ok := gm.face.GlyphBounds(r)
	if !ok {
		return nil, false
	}

	width := int(advance >> 6)
	height := int(bounds.Max.Y - bounds.Min.Y>>6)

	img := image.NewGray(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	baseline := int(-bounds.Min.Y >> 6)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: gm.face,
		Dot:  fixed.Point26_6{X: fixed.I(0), Y: fixed.I(baseline)},
	}

	d.DrawString(string(r))

	meta := &GData{
		char:          r,
		Source:        img,
		DistanceField: generateSDF(img.Pix, height, width),
	}

	return meta, true
}
