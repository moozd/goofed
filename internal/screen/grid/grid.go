package grid

import (
	"image/color"
)

type Grid struct {
	Bg         color.Color
	Size       *GSize
	Cursor     *Cursor
	Cells      []Cell
	viewOffset int
	CellSize   *Size
}

type Cell struct {
	Rune  rune
	Fg    color.Color
	Bg    color.Color
	Dirty bool
}

type Cursor struct {
	Pos    *GPos
	Hidden bool
}

type Size struct {
	Height int
	Width  int
}

type GSize struct {
	Rows int
	Cols int
}

type GPos struct {
	Row int
	Col int
}

type GridViewMode int

const (
	ViewAll GridViewMode = iota
	ViewDirty
)

func New() *Grid {

	cursor := &Cursor{
		Hidden: false,
		Pos: &GPos{
			Col: 0,
			Row: 0,
		}}

	return &Grid{
		Bg:         color.Black,
		Cursor:     cursor,
		viewOffset: 0,
		CellSize:   &Size{Height: 10, Width: 10},
		Size:       &GSize{Cols: 0, Rows: 0},
	}
}

func (self *Grid) Resize(windowWidth, windowHeight int32, blockWidth, blockHeight int32) {
	self.Size.Rows = int(windowWidth) / int(blockWidth)
	self.Size.Cols = int(windowHeight) / int(blockHeight)
	self.CellSize = &Size{Height: int(blockHeight), Width: int(blockWidth)}

	ct := len(self.Cells)
	nt := self.Size.Cols * self.Size.Rows

	if ct < nt {
		self.Cells = append(self.Cells, make([]Cell, nt-ct)...)
		for i := range self.Cells {
			self.Cells[i] = Cell{
				Fg:    color.White,
				Bg:    color.Black,
				Rune:  ' ',
				Dirty: true,
			}

		}
	}

	self.GetView(ViewAll, func(row, col int, cell *Cell) {
		cell.Dirty = true
	})

}

func (self *Grid) getBuffer() []Cell {
	return self.Cells[self.viewOffset*self.Size.Cols : (self.viewOffset*self.Size.Cols)+(self.Size.Cols*self.Size.Rows)]
}

func (self *Grid) GetView(mode GridViewMode, cbl func(x int, y int, cell *Cell)) {
	buff := self.getBuffer()

	for i, cell := range buff {
		if !cell.Dirty && mode == ViewDirty {
			continue
		}

		y := int(i / self.Size.Cols)
		x := i % self.Size.Cols

		cbl(x, y, &cell)
	}
}

func (self *Grid) line(i int) []Cell {
	buff := self.Cells[self.viewOffset*self.Size.Cols:]

	low := i * self.Size.Cols
	high := (i + 1) * self.Size.Cols

	return slice(buff, low+self.viewOffset, high+self.viewOffset)
}

func (self *Grid) getDefaultViewOffset() int {
	return (len(self.Cells) / self.Size.Cols) - self.Size.Rows
}

func (self *Grid) ResetViewOffset() {
	self.viewOffset = self.getDefaultViewOffset()
}

func (self *Grid) Scroll(o int) {
	dvo := self.getDefaultViewOffset()

	if self.Size.Cols == 0 {
		return
	}
	if self.viewOffset+o < 0 {
		self.viewOffset = 0
	} else if self.viewOffset+o > dvo {
		self.viewOffset = dvo
	} else {
		self.viewOffset += o
	}
}
