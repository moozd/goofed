package gfx

import (
	"image"
	"image/draw"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type GMap struct {
	font  *opentype.Font
	face  font.Face
	cache map[rune]*GMeta
}

type GMeta struct {
	Width, Height int
	Source        *image.Gray
	DistanceField *image.Gray
}

func NewGMap(addr string, size int) (*GMap, error) {

	gm := &GMap{
		cache: make(map[rune]*GMeta),
	}

	fb, err := os.ReadFile(addr)
	if err != nil {
		return nil, err
	}

	fnt, err := opentype.Parse(fb)
	if err != nil {
		return nil, err
	}

	fc, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    32,
		DPI:     92,
		Hinting: font.HintingNone,
	})

	gm.face = fc

	return gm, nil
}

func (gm *GMap) Close() {
	gm.face.Close()
}

func (gm *GMap) Get(r rune) (*GMeta, bool) {
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

func (gm *GMap) create(r rune) (*GMeta, bool) {

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

	meta := &GMeta{
		Height:        height,
		Width:         width,
		Source:        img,
		DistanceField: generateSDF(img.Pix, height, width),
	}

	return meta, true
}
