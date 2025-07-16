package font

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"

	"github.com/go-gl/gl/v3.2-core/gl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type GlyphInfo struct {
	X, Y           int
	Width          int
	Height         int
	Advance        float64
	BearingX       float64
	BearingY       float64
	U0, V0, U1, V1 float32
}

type buildOption struct {
	height    int
	width     int
	runeStart rune
	runeEnd   rune
	padding   int

	opentype.FaceOptions
}

type Font struct {
	glyphs map[rune]GlyphInfo
	atlas  image.Image
	font   opentype.Font
	face   font.Face
}

func New(path string, size int) (*Font, error) {

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	self := &Font{}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	ttf, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}

	self.font = *ttf

	self.SetSize(size)
	return self, nil
}
func (self *Font) Measure(r rune) (width, height int) {
	bounds, _, ok := self.face.GlyphBounds(r)
	if !ok {
		return 0, 0 // glyph not available
	}

	width = (bounds.Max.X - bounds.Min.X).Ceil()
	height = (bounds.Max.Y - bounds.Min.Y).Ceil()

	return
}

func (self *Font) Atlas() image.Image {
	return self.atlas
}

func (self *Font) AtlasSize() (width, height int) {
	b := self.atlas.Bounds()
	width = b.Max.X - b.Min.X
	height = b.Max.Y - b.Min.Y
	return

}

func (self *Font) Glyph(r rune) GlyphInfo {
	return self.glyphs[r]
}

func (self *Font) SetSize(s int) {
	err := self.build(&buildOption{
		width:     512,
		height:    512,
		runeStart: rune(32),
		runeEnd:   rune(127), // Change to rune(0x07FF) or higher for Unicode
		padding:   20,
		FaceOptions: opentype.FaceOptions{
			Size:    float64(s),
			DPI:     192,
			Hinting: font.HintingFull,
		},
	})

	if err != nil {
		panic(err)
	}
}

func (self *Font) Bounds() (w int, h int, adv int) {
	metrics := self.face.Metrics()

	bounds, advance, _ := self.face.GlyphBounds('A')

	h = int(metrics.Height)
	w = (bounds.Max.X - bounds.Min.X).Round()
	adv = advance.Round()

	return 0, w, adv

}

func (self *Font) Debug() {
	file, err := os.Create("atlas.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, self.atlas)
	if err != nil {
		panic(err)
	}
	file, err = os.Create("atlas.metadata.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	err = enc.Encode(self.glyphs)
	if err != nil {
		panic(err)
	}
}

func (self *Font) build(options *buildOption) error {

retry:
	face, err := opentype.NewFace(&self.font, &options.FaceOptions)
	if err != nil {
		return err
	}
	defer face.Close()

	self.face = face

	self.glyphs = make(map[rune]GlyphInfo)
	atlas := image.NewRGBA(image.Rect(0, 0, options.width, options.height))
	draw.Draw(atlas, atlas.Bounds(), image.Black, image.Point{}, draw.Src)

	x, y := options.padding, options.padding
	rowHeight := 0

	d := &font.Drawer{
		Dst:  atlas,
		Src:  image.White,
		Face: face,
	}

	metrics := face.Metrics()
	ascent := metrics.Ascent.Ceil()

	for r := options.runeStart; r < options.runeEnd; r++ {
		bounds, advance, _ := face.GlyphBounds(r)
		w := (bounds.Max.X - bounds.Min.X).Ceil()
		h := (bounds.Max.Y - bounds.Min.Y).Ceil()

		if x+w+options.padding > options.width {
			x = options.padding
			y += rowHeight + options.padding
			rowHeight = 0
		}

		if y+h+options.padding > options.height {
			options.width *= 2
			options.height *= 2
			goto retry
		}

		d.Dot = fixed.Point26_6{
			X: fixed.I(x - bounds.Min.X.Ceil()),
			Y: fixed.I(y + ascent),
		}
		d.DrawString(string(r))

		self.glyphs[r] = GlyphInfo{
			X:        x,
			Y:        y,
			Width:    w,
			Height:   h,
			Advance:  float64(advance.Ceil()),
			BearingX: float64(bounds.Min.X.Ceil()),
			BearingY: float64(bounds.Max.Y.Ceil()),
			U0:       float32(x) / float32(options.width),
			V0:       float32(y) / float32(options.height),
			U1:       float32(x+w) / float32(options.width),
			V1:       float32(y+h) / float32(options.height),
		}

		x += w + options.padding
		if h > rowHeight {
			rowHeight = h
		}
	}

	self.atlas = atlas
	return nil

}

func (self *Font) ConfigureTexture(tex uint32) {
	alpha := image.NewAlpha(self.atlas.Bounds())
	draw.Draw(alpha, alpha.Bounds(), self.atlas, image.Point{}, draw.Src)

	gl.BindTexture(gl.TEXTURE_2D, tex)

	w, h := int32(alpha.Rect.Size().X), int32(alpha.Rect.Size().Y)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, w, h, 0,
		gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(alpha.Pix))

	// Swizzle R to Alpha
	swizzle := []int32{gl.ZERO, gl.ZERO, gl.ZERO, gl.RED}
	gl.TexParameteriv(gl.TEXTURE_2D, gl.TEXTURE_SWIZZLE_RGBA, &swizzle[0])

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

}
