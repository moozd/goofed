package gfx

import (
	"fmt"
	"testing"
)

func TestGenerateSDF(t *testing.T) {
	src := [][]byte{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 1, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 0, 0, 0},
	}

	generator := NewSdf(src)

	// add assertions
	debug("src", src)
	debug("dst", generator.buff)

}

func debug[T any](label string, v [][]T) {
	fmt.Printf("%s=[\n", label)
	for _, row := range v {
		fmt.Printf("\t%v\n", row)
	}
	fmt.Println("]")
	fmt.Println()
}
