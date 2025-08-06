package gfx

import "image"

type Atlas struct {
	gm    *GMap
	img   *image.RGBA
	texId uint32
}

// TODO: create a atlas that  has cols/rows like a grid
//   - we ll store uvs and maintain the same atlas in the order that they are being added
//   - Texture width ha limitation on gpu.
//   - the functions(in this case Grid) that will right to screen. needs to add it to Atlas
//     Atlas will store the coords and the rune
//   - during the rendering  we'll call Update so texture will be updated if needed
func NewAtlas(gm *GMap) *Atlas {
	a := &Atlas{gm: gm}

	return a
}

func (a *Atlas) Tex() uint32 {
	return a.texId
}
