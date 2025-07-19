package screen

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

	dirtyCount int
}

type Cell struct {
	Rune  rune
	Fg    color.Color
	Bg    color.Color
	dirty bool
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

type ViewMode int

const (
	GridIterAll ViewMode = iota
	GridIterDirty
)

func newGrid() *Grid {

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

	self.Size.Rows = int(windowHeight) / int(blockHeight)
	self.Size.Cols = int(windowWidth) / int(blockWidth)
	self.CellSize = &Size{Height: int(blockHeight), Width: int(blockWidth)}

	ct := len(self.Cells)
	nt := self.Size.Cols * self.Size.Rows

	if ct < nt {
		self.Cells = append(self.Cells, make([]Cell, nt-ct)...)
		for i := range self.Cells {
			self.Cells[i] = Cell{
				Bg:   color.White,
				Fg:   color.Black,
				Rune: 'A',
			}
		}
	}

	self.GetView(GridIterAll, func(row, col int, cell *Cell) {
		self.markDirty(cell)
	})

}

func (self *Grid) IsClean() bool {
	return self.dirtyCount == 0

}

func (self *Grid) markDirty(cell *Cell) {
	self.dirtyCount = +1
	cell.dirty = true
}
func (self *Grid) markClean(cell *Cell) {
	self.dirtyCount -= 1
	if self.dirtyCount < 0 {
		self.dirtyCount = 0
	}
	cell.dirty = false
}

func (self *Grid) getBuffer() []Cell {
	return self.Cells[self.viewOffset*self.Size.Cols : (self.viewOffset*self.Size.Cols)+(self.Size.Cols*self.Size.Rows)]
}

func (self *Grid) GetView(mode ViewMode, cbl func(x int, y int, cell *Cell)) {
	buff := self.getBuffer()

	for i := range buff {
		cell := &buff[i]
		if !cell.dirty && mode == GridIterDirty {
			continue
		}

		y := int(i / self.Size.Cols)
		x := i % self.Size.Cols

		cbl(x, y, cell)

		if mode == GridIterDirty {
			self.markClean(cell)
		}
	}
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
