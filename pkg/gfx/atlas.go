package gfx

import (
	"image"
	"image/draw"
	"log"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Atlas struct {
	gm         *Font
	img        *image.RGBA
	texId      uint32
	dirty      bool
	queue      []*Glyph
	meta       map[rune]*atlasGMeta
	row, col   int
	rows, cols int
	total      int
}

type atlasGMeta struct {
	added bool
	X, Y  int
}

func NewAtlas(gm *Font) *Atlas {
	a := &Atlas{gm: gm}

	a.dirty = true
	a.initAtlasSize(256)

	a.meta = make(map[rune]*atlasGMeta)
	a.img = image.NewRGBA(image.Rect(0, 0, a.cols*gm.AdvanceWidth, a.rows*gm.LineHeight))
	draw.Draw(a.img, a.img.Bounds(), image.Black, image.Point{}, draw.Src)

	return a
}

func (a *Atlas) GetSize() (w, h float32) {
	bounds := a.img.Bounds()
	w = float32(bounds.Dx())
	h = float32(bounds.Dy())
	return
}

func (a *Atlas) initAtlasSize(rectWidth int) {
	cw := a.gm.AdvanceWidth

	a.cols = int(math.Ceil(float64(rectWidth) / float64(cw)))
	a.rows = a.cols
}

func (a *Atlas) advance() {
	a.col = (a.total - 1) % a.cols
	a.row = int(math.Floor(float64((a.total - 1) / a.cols)))
	log.Printf("row: %d, col: %d", a.row, a.col)
}

func (a *Atlas) Update(runes ...rune) {

	for _, c := range runes {
		if _, ok := a.meta[c]; ok {
			continue
		}

		data, ok := a.gm.Get(c)
		if !ok {
			continue
		}

		a.total += 1

		a.meta[c] = &atlasGMeta{
			added: true,
			X:     a.col * a.gm.AdvanceWidth,
			Y:     a.row * a.gm.LineHeight,
		}
		a.advance()

		log.Printf("X:%d ,Y:%d\n", a.meta[c].X, a.meta[c].Y)

		a.queue = append(a.queue, data)

		a.dirty = true
	}

}

func (a *Atlas) GetUVs(r rune) (u0, v0, u1, v1 float32) {
	m := a.meta[r]
	W := float32(a.img.Bounds().Dx())
	H := float32(a.img.Bounds().Dy())

	u0 = float32(m.X) / W
	u1 = float32(m.X+a.gm.AdvanceWidth) / W
	v0 = float32(m.Y) / H
	v1 = float32(m.Y+a.gm.LineHeight) / H
	return
}

func (a *Atlas) Compile() {
	if !a.dirty {
		return
	}

	for len(a.queue) > 0 {

		tex := a.queue[0]
		a.queue = a.queue[1:]

		src := tex.DistanceField

		m := a.meta[tex.char]

		x, y := m.X, m.Y

		draw.Draw(a.img,
			image.Rect(x, y, x+src.Bounds().Dx(), y+src.Bounds().Dy()),
			src,
			image.Point{0, 0},
			draw.Over)

	}

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	diagnose()
	gl.GenTextures(1, &a.texId)
	diagnose()
	gl.ActiveTexture(gl.TEXTURE0)
	diagnose()
	gl.BindTexture(gl.TEXTURE_2D, a.texId)
	diagnose()
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(a.img.Rect.Size().X),
		int32(a.img.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(a.img.Pix),
	)
	diagnose()

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	diagnose()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	diagnose()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	diagnose()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	diagnose()

	a.dirty = false

}

func (a *Atlas) TexSlotIndex() int32 {
	return 0
}

func (t *Atlas) Delete() {
	gl.DeleteTextures(1, &t.texId)
}
