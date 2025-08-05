package gfx

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func TestGenerateSDF(t *testing.T) {

	fontBytes, _ := os.ReadFile("/home/mo/.local/share/fonts/FiraCode/FiraCodeNerdFont-Regular.ttf")

	// Parse the font
	f, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal("Failed to parse font:", err)
	}

	// Get the font's units per em (original design size)
	unitsPerEm := f.UnitsPerEm()

	// Create font face at original size (units per em)
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(unitsPerEm) / 4,
		DPI:     12,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal("Failed to create font face:", err)
	}
	defer face.Close()

	// Get metrics for the glyph
	glyph := 'A' // The glyph to draw
	bounds, advance, ok := face.GlyphBounds(glyph)
	if !ok {
		log.Fatal("Glyph not found in font")
	}

	// Calculate image dimensions based on glyph bounds
	width := int((advance >> 6) + 20)                  // Add some padding
	height := int((bounds.Max.Y-bounds.Min.Y)>>6) + 20 // Add some padding

	// Ensure minimum size
	if width < 50 {
		width = 50
	}
	if height < 50 {
		height = 50
	}

	// Create a grayscale image sized to fit the glyph
	img := image.NewGray(image.Rect(0, 0, width, height))

	// Fill with white background
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	// Calculate baseline position
	baseline := int(-bounds.Min.Y>>6) + 10 // 10 pixels padding from top

	// Create a drawer for text
	d := &font.Drawer{
		Dst:  img,
		Src:  image.White, // Black text
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(10), Y: fixed.I(baseline)}, // Position with padding
	}

	// Draw the glyph at its original size
	d.DrawString(string(glyph))
	src := img.Pix

	// Save the original binary image for comparison
	z, _ := os.Create("A_original.png")
	png.Encode(z, img)
	z.Close()

	sdf := NewSdf(src, height, width)

	sdf.Save("A.png")

}

func debug[T any](label string, v [][]T) {
	fmt.Printf("%s=[\n", label)
	for _, row := range v {
		fmt.Printf("\t%v\n", row)
	}
	fmt.Println("]")
	fmt.Println()
}
