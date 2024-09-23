package cellular

import (
	"fmt"
	"math"
	"time"
)

// canvas represents a grid of virtual pixels (i.e. cells) for our automata, as it is unlikely we want every single real pixel to be
type canvas struct {
	cells [][]cell
	// cell counts
	width, height uint
	// real pixel counts
	realWidth, realHeight uint
	// cell dimensions in pixels
	cellWidth, cellHeight uint
}

func (c canvas) at(x, y int) (cell, error) {
	if x < 0 || x >= len(c.cells) {
		return 0, fmt.Errorf("index x=%v off grid", x)
	}

	if y < 0 || y >= len(c.cells[0]) {
		return 0, fmt.Errorf("index y=%v off grid", y)
	}

	return c.cells[x][y], nil
}

func newCanvas(width, height, realWidth, realHeight uint) canvas {
	c := canvas{}
	c.cells = make([][]cell, width)
	for x := range c.cells {
		c.cells[x] = make([]cell, height)
	}
	c.width = width
	c.height = height
	c.realWidth = realWidth
	c.realHeight = realHeight

	c.cellWidth = c.realWidth / c.width
	c.cellHeight = c.realHeight / c.height

	return c
}

// cell represents the state of a single cell, as defined in cellStateModel
type cell uint

type predicate func(canvas, uint, uint) bool

type transition struct {
	predicate predicate
	newState  cell
}

type rgb struct {
	r, g, b float64
}

type cellColouring []rgb

type cellularAutomata interface {
	getCellColouring() cellColouring
	getTransitionModel() [][]transition
}

func Render(drawChan chan<- []float64, doneChan <-chan any, fps, width, height uint, pixelsX, pixelsY uint) {
	// this must be run as a goroutine
	go func(drawChan chan<- []float64, doneChan <-chan any, fps, width, height uint, pixelsX, pixelsY uint) {
		// barf shit down drawchan to render it

		// convert fps into seconds per frame
		frameDuration := time.Duration(math.Pow(float64(fps), -1)*1000) * time.Millisecond
		fpsClock := time.NewTicker(frameDuration)

		c := newCanvas(width, height, pixelsX, pixelsY)

		c.cells[0][2] = 1
		c.cells[1][0] = 1
		c.cells[1][2] = 1
		c.cells[2][1] = 1
		c.cells[2][2] = 1

		conways := newConways()

		for {
			select {
			case <-fpsClock.C:
				c = step(c, conways)
				pixelGrid := paint(c, conways)
				drawChan <- pixelGrid
			case <-doneChan:
				close(drawChan)
				return
			}
		}
	}(drawChan, doneChan, fps, width, height, pixelsX, pixelsY)
}

func step(c canvas, automata cellularAutomata) canvas {
	new := newCanvas(c.width, c.height, c.realWidth, c.realHeight)

	model := automata.getTransitionModel()

	for x := range len(c.cells) {
		for y := range len(c.cells[0]) {
			thisCell, err := c.at(x, y)
			if err != nil {
				panic("Something has gone very wrong. We indexed outside of the grid in the step function.")
			}
			for _, t := range model[thisCell] {
				if t.predicate(c, uint(x), uint(y)) {
					new.cells[x][y] = t.newState
					break
				}
			}
		}
	}

	return new
}

func paint(c canvas, automata cellularAutomata) []float64 {
	colourings := automata.getCellColouring()

	pixels := make([]float64, 4*c.realWidth*c.realHeight)

	for x := range c.width {
		for y := range c.height {
			thisCell, err := c.at(int(x), int(y))
			if err != nil {
				panic("Something has gone very wrong. We indexed outside of the grid in the step function.")
			}
			setPixel(pixels, uint(x), uint(y), colourings[thisCell], c)
		}
	}

	return pixels
}

func setPixel(pixels []float64, x, y uint, colour rgb, c canvas) []float64 {
	pixelLeftBound := c.cellWidth * x
	pixelRightBound := pixelLeftBound + c.cellWidth - 1
	pixelLowerBound := c.cellHeight * y
	pixelUpperBound := pixelLowerBound + c.cellHeight - 1

	for thisX := pixelLeftBound; thisX <= pixelRightBound; thisX++ {
		for thisY := pixelLowerBound; thisY <= pixelUpperBound; thisY++ {
			redIndex := getRealPixelIndex(thisX, thisY, c.realWidth)
			pixels[redIndex] = colour.r
			pixels[redIndex+1] = colour.g
			pixels[redIndex+2] = colour.b
			pixels[redIndex+3] = 1
		}
	}

	return pixels
}

func getRealPixelIndex(realX, realY, canvasWidth uint) uint {
	return 4 * (realY*canvasWidth + realX)
}

type conways struct {
	model     [][]transition
	colouring cellColouring
}

func (c conways) getTransitionModel() [][]transition {
	return c.model
}

func (c conways) getCellColouring() cellColouring {
	return c.colouring
}

func newConways() cellularAutomata {
	conways := conways{}

	model := make([][]transition, 0)

	countAliveNeighbours := func(c canvas, x, y int) uint {
		neighbours := uint(0)
		for offsetX := -1; offsetX <= 1; offsetX++ {
			for offsetY := -1; offsetY <= 1; offsetY++ {
				if offsetY == 0 && offsetX == 0 {
					continue
				}

				neighbour, err := c.at(x+offsetX, y+offsetY)
				if err != nil {
					// gone off the grid, there's no neighbour here
					continue
				}

				if neighbour == 1 {
					neighbours++
				}
			}
		}
		return neighbours
	}

	// dead
	state0 := make([]transition, 0)
	state0 = append(state0, transition{
		// 2 or 3 alive neighbours, become alive
		predicate: func(c canvas, x, y uint) bool {
			return countAliveNeighbours(c, int(x), int(y)) == 3
		},
		newState: 1,
	})
	model = append(model, state0)

	state1 := make([]transition, 0)
	state1 = append(state1, transition{
		// <2 alive neighbours, die
		predicate: func(c canvas, x, y uint) bool {
			return countAliveNeighbours(c, int(x), int(y)) < 2
		},
		newState: 0,
	})
	state1 = append(state1, transition{
		// >3 alive neighbours, die
		predicate: func(c canvas, x, y uint) bool {
			return countAliveNeighbours(c, int(x), int(y)) > 3
		},
		newState: 0,
	})
	state1 = append(state1, transition{
		// 2 or 3 alive neighbours, live
		predicate: func(c canvas, x, y uint) bool {
			alive := countAliveNeighbours(c, int(x), int(y))
			return alive == 2 || alive == 3
		},
		newState: 1,
	})
	model = append(model, state1)

	conways.model = model
	conways.colouring = make(cellColouring, 2)
	conways.colouring[0] = rgb{0, 0, 0}
	conways.colouring[1] = rgb{1, 1, 1}

	return conways
}
