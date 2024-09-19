package cellular

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// canvas represents a grid of virtual pixels (i.e. cells) for our automata, as it is unlikely we want every single real pixel to be
type canvas struct {
	cells [][]cell
	// represents real pixel counts
	width, height uint
}

func (c canvas) at(x, y uint) cell {
	return c.cells[x][y]
}

func newCanvas(width, height, realWidth, realHeight uint) canvas {
	c := canvas{}
	c.cells = make([][]cell, width)
	for x := range c.cells {
		c.cells[x] = make([]cell, height)
	}
	c.width = realWidth
	c.height = realHeight
	return c
}

// cell represents the state of a single cell, as defined in cellStateModel
type cell uint

type predicate func(canvas, uint, uint) bool

type transition struct {
	p        predicate
	newState cell
}

type rgb struct {
	r, g, b float64
}

type cellColouring []rgb

type cellularAutomata interface {
	getCellColouring() cellColouring
	getTransitionModel() [][]transition
}

/* thinking

I want cell states to be a natural number. For example, in conway's game of life, we model 0 as dead and 1 as alive

to model state transitions, we can have an array that is indexed by the current state, say 1 (alive)

the elements of the array needs to be some sort of list of conditions. the first one found that is true, changes the cell's current state

for example, in conway's:

[
	{ // index 0
		are all neighbouring cells dead? -> 0
		are all neighbouring cells alive? -> 1
	},
	{ // index 1
		are there two or three neighbouring cells alive? -> 1
		otherwise -> 0
	}
]

(something like that, anyway)

So we need a way of encoding these questions as a predicate

The only sane way I can think of is a func(canvas) bool that is provided by config or something, that queries our canvas

whoops, map doesn't make sense. What about an array of predicates? or a map from state to []predicate ???
*/

func Render(drawChan chan<- []float64, doneChan <-chan any, fps, height, width uint, pixelsY, pixelsX uint) {
	// this must be run as a goroutine
	go func(drawChan chan<- []float64, doneChan <-chan any, fps, height, width uint) {
		// barf shit down drawchan to render it

		// fps^-1 * 1000 = milliseconds per frame
		// 10fps -> (10)^-1 * 1000 = 0.1 * 1000 = 100
		frameDuration := time.Duration(math.Pow(float64(fps), -1)*1000) * time.Millisecond
		fpsClock := time.NewTicker(frameDuration)

		c := newCanvas(width, height)
		c.cells[3][2] = 1
		c.cells[3][3] = 1
		c.cells[3][4] = 1

		conways := newConways()

		pixelGrid := make([]float64, 4*pixelsX*pixelsY)

		for {
			select {
			case <-fpsClock.C:
				c = step(c, conways)
				err := paint(&pixelGrid, c, conways)
				if err != nil {
					close(drawChan)
					panic(fmt.Errorf("error painting grid: %w", err))
				}
				drawChan <- pixelGrid
			case <-doneChan:
				close(drawChan)
				return
			}
		}
	}(drawChan, doneChan, fps, height, width)
}

func step(c canvas, automata cellularAutomata) canvas {
	new := newCanvas(uint(len(c.cells)), uint(len(c.cells[0])))

	model := automata.getTransitionModel()

	for x := range len(c.cells) {
		for y := range len(c.cells[0]) {
			for _, t := range model[c.at(uint(x), uint(y))] {
				if t.p(c, uint(x), uint(y)) {
					new.cells[x][y] = t.newState
					break
				}
			}
		}
	}

	return new
}

func paint(pixels *[]float64, c canvas, automata cellularAutomata) error {
	if len(*pixels)%4 != 0 {
		return errors.New("supplied pixel array length not a multiple of 4")
	}

	colourings := automata.getCellColouring()

	virtualWidth := len(c.cells)
	virtualHeight := len(c.cells[0])

	pixelWidth := c.width / uint(virtualWidth)
	pixelHeight := c.height / uint(virtualHeight)

	i := 0
	for x := range virtualWidth {
		for y := range virtualHeight {
			colour := colourings[c.at(uint(x), uint(y))]
			(*pixels)[i] = colour.r
			(*pixels)[i+1] = colour.g
			(*pixels)[i+2] = colour.b
			(*pixels)[i+3] = 1
			i += 4
		}
	}

	return nil
}

func setPixel(pixels *[]float64, x, y uint, colour rgb) {
	// urgh, some hard maths here to work out which elements to modify. I'm doing this tomorrow, laters
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

	// dead
	state0 := make([]transition, 0)
	state0 = append(state0, transition{
		// 3 alive neighbours, become alive
		p: func(c canvas, x, y uint) bool {
			// todo
			return false
		},
		newState: 1,
	})
	model = append(model, state0)

	state1 := make([]transition, 0)
	state1 = append(state1, transition{
		// <2 alive neighbours, die
		p: func(c canvas, x, y uint) bool {
			// todo
			return false
		},
		newState: 0,
	})
	state1 = append(state1, transition{
		// >3 alive neighbours, die
		p: func(c canvas, x, y uint) bool {
			// todo
			return false
		},
		newState: 0,
	})
	model = append(model, state1)

	conways.model = model
	conways.colouring = make(cellColouring, 2)
	conways.colouring[0] = rgb{0, 0, 0}
	conways.colouring[1] = rgb{1, 1, 1}

	return conways
}
