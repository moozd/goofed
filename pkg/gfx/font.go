package gfx

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	texId              uint32
	cw, ch, cols, rows int
}

func (t *Font) TexID() uint32 {
	return t.texId
}

func NewFont(path string, size float64) *Font {
	t := &Font{}
	startChar := 32
	endChar := 128
	t.cols = 16
	t.rows = 6
	t.ch = int(size)
	t.cw = int(size)

	fb := assert(os.ReadFile(path))
	ft := assert(opentype.Parse(fb))
	fc := assert(opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}))

	aw := t.cols * t.cw
	ah := t.rows * t.ch
	img := image.NewRGBA(image.Rect(0, 0, aw, ah))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: fc,
	}

	metrics := fc.Metrics()
	ascent := metrics.Ascent.Round()
	// descent := metrics.Descent.Round()
	// lineHeight := ascent + descent

	charIndex := 0
	for c := startChar; c <= endChar; c++ {
		x := (charIndex % t.cols) * t.cw
		y := (charIndex / t.cols) * t.ch

		// Optional: debug background
		cellRect := image.Rect(x, y, x+t.cw, y+t.ch)
		draw.Draw(img, cellRect, &image.Uniform{C: color.RGBA{30, 30, 30, 255}}, image.Point{}, draw.Src)

		// Correct baseline
		d.Dot = fixed.P(x, y+ascent)
		d.DrawString(string(rune(c)))

		charIndex++
	}

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	diagnose()
	gl.GenTextures(1, &t.texId)
	diagnose()
	gl.ActiveTexture(gl.TEXTURE0)
	diagnose()
	gl.BindTexture(gl.TEXTURE_2D, t.texId)
	diagnose()
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(img.Rect.Size().X),
		int32(img.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pix),
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

	outFile := assert(os.Create("font_atlas.png"))
	defer outFile.Close()

	assert(0, png.Encode(outFile, img))

	return t
}

func (f *Font) GetUVs(r rune) (u0, v0, u1, v1 float32) {
	charIndex := int(r) - 32
	col := charIndex % f.cols
	row := charIndex / f.cols

	// Flip the row to match OpenGL UV origin (bottom-left)
	flippedRow := f.rows - row - 1

	charWUV := 1.0 / float32(f.cols)
	charHUV := 1.0 / float32(f.rows)

	u0 = float32(col) * charWUV
	u1 = float32(col+1) * charWUV
	v0 = float32(flippedRow) * charHUV
	v1 = float32(flippedRow+1) * charHUV

	log.Printf("Font %c  u0=%f v0=%f u1=%f v11=%f", u0, v0, u1, v1)
	return
}

func (t *Font) Delete() {
	gl.DeleteTextures(1, &t.texId)
}
