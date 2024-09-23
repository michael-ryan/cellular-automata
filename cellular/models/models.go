package models

import "fmt"

// Canvas represents a grid of virtual pixels (i.e. cells) for our automata, as it is unlikely we want every single real pixel to be simulated as a cell
type Canvas struct { // todo this really shouldnt live here. also other things. move shit around to make sense
	Cells [][]Cell
	// cell counts
	Width, Height uint
	// real pixel counts
	RealWidth, RealHeight uint
	// cell dimensions in pixels
	CellWidth, CellHeight uint
}

func (c Canvas) At(x, y int) (Cell, error) {
	if x < 0 || x >= len(c.Cells) {
		return 0, fmt.Errorf("index x=%v off grid", x)
	}

	if y < 0 || y >= len(c.Cells[0]) {
		return 0, fmt.Errorf("index y=%v off grid", y)
	}

	return c.Cells[x][y], nil
}

func NewCanvas(width, height, realWidth, realHeight uint) Canvas {
	c := Canvas{}
	c.Cells = make([][]Cell, width)
	for x := range c.Cells {
		c.Cells[x] = make([]Cell, height)
	}
	c.Width = width
	c.Height = height
	c.RealWidth = realWidth
	c.RealHeight = realHeight

	c.CellWidth = c.RealWidth / c.Width
	c.CellHeight = c.RealHeight / c.Height

	return c
}

// Cell represents the state of a single Cell, as defined in cellStateModel
type Cell uint

type Predicate func(Canvas, uint, uint) bool

type Transition struct {
	Predicate Predicate
	NewState  Cell
}

type Rgb struct {
	R, G, B float64
}

type cellColouring []Rgb

type CellularAutomata interface {
	GetCellColouring() cellColouring
	GetTransitionModel() [][]Transition
}
