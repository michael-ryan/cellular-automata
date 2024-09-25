package cellular

import (
	"math"
	"sync"
	"time"

	"github.com/michael-ryan/cellular-automata/cellular/models"
)

type Coord [2]uint

func StartRenderer(drawChan chan<- []float64, startChan <-chan any, clickChan <-chan Coord, doneChan <-chan any, fps, width, height uint, pixelsX, pixelsY uint) {
	// this must be run as a goroutine
	go func(drawChan chan<- []float64, doneChan <-chan any, fps, width, height uint, pixelsX, pixelsY uint) {
		// barf pixel arrays down drawchan to render it

		// convert fps into seconds per frame
		frameDuration := time.Duration(math.Pow(float64(fps), -1)*1000) * time.Millisecond
		fpsClock := time.NewTicker(frameDuration)

		c := models.NewCanvas(width, height, pixelsX, pixelsY)

		automata := models.NewLangtons()

		started := false

		for !started {
			select {
			case <-fpsClock.C:
				pixelGrid := paint(c, automata)
				drawChan <- pixelGrid
			case click, ok := <-clickChan:
				if !ok {
					continue
				}
				location := getVirtualPixelXY(click, c)
				oldCell := c.Cells[location[0]][location[1]]
				newCell := (oldCell + 1) % models.Cell(len(automata.GetTransitionModel()))
				c.Cells[location[0]][location[1]] = newCell
			case <-startChan:
				started = true
			}
		}

		for {
			select {
			case <-fpsClock.C:
				c = step(c, automata)
				pixelGrid := paint(c, automata)
				drawChan <- pixelGrid
			case <-doneChan:
				close(drawChan)
				return
			}
		}
	}(drawChan, doneChan, fps, width, height, pixelsX, pixelsY)
}

func step(c models.Canvas, automata models.CellularAutomata) models.Canvas {
	new := models.NewCanvas(c.Width, c.Height, c.RealWidth, c.RealHeight)
	new.Cells = make([][]models.Cell, len(c.Cells))
	for x := range len(new.Cells) {
		new.Cells[x] = make([]models.Cell, len(c.Cells[0]))
		for y := range len(new.Cells[x]) {
			new.Cells[x][y] = c.Cells[x][y]
		}
	}

	model := automata.GetTransitionModel()

	type edit struct {
		x, y     int
		newState models.Cell
	}

	wg := sync.WaitGroup{}
	editChan := make(chan edit)

	for x := range len(c.Cells) {
		for y := range len(c.Cells[0]) {
			// run each cell compute in its own goroutine
			wg.Add(1)
			go func(x, y int, c models.Canvas, model [][]models.Transition, editChan chan<- edit) {
				defer wg.Done()
				thisCell, err := c.At(x, y)
				if err != nil {
					panic("Something has gone very wrong. We indexed outside of the grid in the step function.")
				}
				for _, t := range model[thisCell] {
					if t.Predicate(c, x, y) {
						editChan <- edit{
							x:        x,
							y:        y,
							newState: t.NewState,
						}
						return
					}
				}
			}(x, y, c, model, editChan)
		}
	}

	go func() {
		wg.Wait()
		close(editChan)
	}()

	for edit := range editChan {
		new.Cells[edit.x][edit.y] = edit.newState
	}

	automata.Step()

	return new
}

func paint(c models.Canvas, automata models.CellularAutomata) []float64 {
	colourings := automata.GetCellColouring()

	pixels := make([]float64, 4*c.RealWidth*c.RealHeight)

	for x := range c.Width {
		for y := range c.Height {
			thisCell, err := c.At(int(x), int(y))
			if err != nil {
				panic("Something has gone very wrong. We indexed outside of the grid in the step function.")
			}
			pixels = setPixel(pixels, uint(x), uint(y), colourings[thisCell], c)
		}
	}

	return pixels
}

func setPixel(pixels []float64, x, y uint, colour models.Rgb, c models.Canvas) []float64 {
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

func getVirtualPixelXY(coord Coord, c models.Canvas) Coord {
	x := coord[0] / c.CellWidth
	y := coord[1] / c.CellHeight
	return Coord{x, y}
}

func getRealPixelIndex(realX, realY, canvasWidth uint) uint {
	return 4 * (realY*canvasWidth + realX)
}
