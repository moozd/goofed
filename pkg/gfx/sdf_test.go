package gfx

import (
	"image"
	"image/png"
	"os"
	"testing"
)

func TestGenerateSDF(t *testing.T) {

	gm, _ := NewGMap("/home/mo/.local/share/fonts/FiraCode/FiraCodeNerdFont-Regular.ttf", 18)
	atlas := NewAtlas(gm)

	atlas.Update('a', 'b', 'c', 'd', 'e', 'f', '1', '2', '3', '@')
	atlas.GetTexID()

	save("atlas.png", atlas.img)

}

func save(name string, img image.Image) {
	file1, _ := os.Create(name)
	defer file1.Close()
	png.Encode(file1, img)
}
