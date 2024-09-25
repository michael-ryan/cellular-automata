package render

import (
	"math"
	"time"

	"github.com/michael-ryan/cellular-automata/render/models"
)

type Coord [2]uint

// canvas represents a grid of virtual pixels (i.e. cells) for our automata, as it is unlikely we want every single real pixel to be simulated as a cell
type canvas struct { // todo this really shouldnt live here. also other things. move stuff around to make sense
	Cells [][]uint
	// cell counts
	Width, Height uint
	// real pixel counts
	RealWidth, RealHeight uint
	// cell dimensions in pixels
	CellWidth, CellHeight uint
}

func NewCanvas(width, height, realWidth, realHeight uint) canvas {
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

func StartRenderer(drawChan chan<- []float64, startChan <-chan any, clickChan <-chan Coord, doneChan <-chan any, fps, width, height uint, pixelsX, pixelsY uint) {
	// this must be run as a goroutine
	go func(drawChan chan<- []float64, doneChan <-chan any, fps, width, height uint, pixelsX, pixelsY uint) {
		// barf pixel arrays down drawchan to render it

		// convert fps into seconds per frame
		frameDuration := time.Duration(math.Pow(float64(fps), -1)*1000) * time.Millisecond
		fpsClock := time.NewTicker(frameDuration)

		c := NewCanvas(width, height, pixelsX, pixelsY)

		automata := models.NewLangtons() // todo parameterise me

		started := false

		for !started {
			select {
			case <-fpsClock.C:
				pixelGrid := paint(c, automata.GetColouring())
				drawChan <- pixelGrid
			case click, ok := <-clickChan:
				if !ok {
					continue
				}
				location := getVirtualPixelXY(click, c)
				oldCell := c.Cells[location[0]][location[1]]
				newCell := (oldCell + 1) % automata.CountStates()
				c.Cells[location[0]][location[1]] = newCell
			case <-startChan:
				started = true
			}
		}

		for {
			select {
			case <-fpsClock.C:
				c.Cells = automata.Step(c.Cells)
				pixelGrid := paint(c, automata.GetColouring())
				drawChan <- pixelGrid
			case <-doneChan:
				close(drawChan)
				return
			}
		}
	}(drawChan, doneChan, fps, width, height, pixelsX, pixelsY)
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

func getVirtualPixelXY(coord Coord, c canvas) Coord {
	x := coord[0] / c.CellWidth
	y := coord[1] / c.CellHeight
	return Coord{x, y}
}

func getRealPixelIndex(realX, realY, canvasWidth uint) uint {
	return 4 * (realY*canvasWidth + realX)
}
