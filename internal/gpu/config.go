package gpu

import "github.com/moozd/goofed/internal/gpu/font"

type Config struct {
	WindowTitle  string
	WindowHeight int
	WindowWidth  int
	Font         *font.Font
}
