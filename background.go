package render

import (
	"math"
	"time"

	"github.com/michael-ryan/cellular-automata/models"
)

type coord [2]uint

// canvas represents a grid of virtual pixels (i.e. cells) for our simulation, as it is unlikely we want every single real pixel to be simulated as a cell
type canvas struct { // todo this really shouldnt live here. also other things. move stuff around to make sense
	Cells [][]uint
	// cell counts
	Width, Height uint
	// real pixel counts
	RealWidth, RealHeight uint
	// cell dimensions in pixels
	CellWidth, CellHeight uint
}

func newCanvas(width, height, realWidth, realHeight uint) canvas {
	c := canvas{}
	c.Cells = make([][]uint, width)
	for x := range c.Cells {
		c.Cells[x] = make([]uint, height)
	}
	c.Width = width
	c.Height = height
	c.RealWidth = realWidth
	c.RealHeight = realHeight

	c.CellWidth = c.RealWidth / c.Width
	c.CellHeight = c.RealHeight / c.Height

	return c
}

func startRenderer(drawChan chan<- []float64, startChan <-chan any, clickChan <-chan coord, doneChan <-chan any, c Config) {
	// this must be run as a goroutine
	go func(drawChan chan<- []float64, doneChan <-chan any, c Config) {
		// convert fps into seconds per frame
		frameDuration := time.Duration(math.Pow(float64(c.Fps), -1)*1000) * time.Millisecond
		fpsClock := time.NewTicker(frameDuration)

		canvas := newCanvas(c.CellsX, c.CellsY, c.WindowX, c.WindowY)

		started := false

		for !started {
			select {
			case <-fpsClock.C:
				pixelGrid := paint(canvas, c.Automaton.GetColouring())
				drawChan <- pixelGrid
			case click, ok := <-clickChan:
				if !ok {
					continue
				}
				location := getVirtualPixelXY(click, canvas)
				oldCell := canvas.Cells[location[0]][location[1]]
				newCell := (oldCell + 1) % c.Automaton.CountStates()
				canvas.Cells[location[0]][location[1]] = newCell
			case <-startChan:
				started = true
			}
		}

		for {
			select {
			case <-fpsClock.C:
				canvas.Cells = c.Automaton.Step(canvas.Cells)
				pixelGrid := paint(canvas, c.Automaton.GetColouring())
				drawChan <- pixelGrid
			case <-doneChan:
				close(drawChan)
				return
			}
		}
	}(drawChan, doneChan, c)
}

func paint(c canvas, colourings []models.Rgb) []float64 {
	pixels := make([]float64, 4*c.RealWidth*c.RealHeight)

	for x := range c.Width {
		for y := range c.Height {
			thisCell, err := models.At(c.Cells, int(x), int(y))
			if err != nil {
				panic("Something has gone very wrong. We indexed outside of the grid in the step function.")
			}
			pixels = setPixel(pixels, uint(x), uint(y), colourings[thisCell], c)
		}
	}

	return pixels
}

func setPixel(pixels []float64, x, y uint, colour models.Rgb, c canvas) []float64 {
	pixelLeftBound := c.CellWidth * x
	pixelRightBound := pixelLeftBound + c.CellWidth - 1
	pixelLowerBound := c.CellHeight * y
	pixelUpperBound := pixelLowerBound + c.CellHeight - 1

	for thisX := pixelLeftBound; thisX <= pixelRightBound; thisX++ {
		for thisY := pixelLowerBound; thisY <= pixelUpperBound; thisY++ {
			redIndex := getRealPixelIndex(thisX, thisY, c.RealWidth)
			pixels[redIndex] = colour.R
			pixels[redIndex+1] = colour.G
			pixels[redIndex+2] = colour.B
			pixels[redIndex+3] = 1
		}
	}

	return pixels
}

func getVirtualPixelXY(xy coord, c canvas) coord {
	x := xy[0] / c.CellWidth
	y := xy[1] / c.CellHeight
	return coord{x, y}
}

func getRealPixelIndex(realX, realY, canvasWidth uint) uint {
	return 4 * (realY*canvasWidth + realX)
}
