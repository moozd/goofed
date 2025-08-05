package gfx

import (
	"fmt"
	"image/png"
	"os"
	"testing"
)

func TestGenerateSDF(t *testing.T) {

	gm, _ := NewGMap("/home/mo/.local/share/fonts/FiraCode/FiraCodeNerdFont-Regular.ttf", 14)

	m, ok := gm.Get('A')
	if ok {
		SaveGMeta('A', m)
	}
	m, ok = gm.Get('B')
	if ok {
		SaveGMeta('A', m)
	}

}

func SaveGMeta(c rune, m *GMeta) {

	file1, _ := os.Create(fmt.Sprintf("%c_og.png", c))
	defer file1.Close()
	png.Encode(file1, m.Image)

	file2, _ := os.Create(fmt.Sprintf("%c_tex.png", c))
	defer file2.Close()
	png.Encode(file2, m.Tex)

}
