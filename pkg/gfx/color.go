package gfx

import "image/color"

type OpenGLColor struct {
	R, G, B, A float32
}

func toOpenGLColor(c color.RGBA) (r, g, b, a float32) {
	r = float32(c.R / 255)
	g = float32(c.G / 255)
	b = float32(c.B / 255)
	a = float32(c.A)
	return
}
