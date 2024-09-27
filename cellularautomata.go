package cellularautomata

import (
	"fmt"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/michael-ryan/cellularautomata/models"
)

type Config struct {
	Fps, CellsX, CellsY, WindowX, WindowY uint
	Automaton                             *models.Automaton
}

// canvas represents a grid of virtual pixels (i.e. cells) for our simulation, as it is unlikely we want every single real pixel to be simulated as a cell
type canvas struct {
	Cells [][]uint
	// cell counts
	Width, Height uint
	// real pixel counts
	RealWidth, RealHeight uint
	// cell dimensions in pixels
	CellWidth, CellHeight uint
}

var config Config

// Launch opens a GUI window that renders the simulation.
// This function will block until the window is closed.
func Launch(c Config) error {
	if c.CellsX > c.WindowX {
		return fmt.Errorf("cellsX (%v) cannot be larger than windowX (%v), since each cell requires at least one pixel", c.CellsX, c.WindowX)
	}

	if c.CellsX > c.WindowX {
		return fmt.Errorf("cellsY (%v) cannot be larger than windowY (%v), since each cell requires at least one pixel", c.CellsY, c.WindowY)
	}

	config = c

	opengl.Run(launch)

	return nil
}

func launch() {
	cfg := opengl.WindowConfig{
		Title:  "Cellular Automata",
		Bounds: pixel.R(0, 0, float64(config.WindowX), float64(config.WindowY)),
		VSync:  true,
	}

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	frameDuration := time.Second / time.Duration(config.Fps)
	fpsClock := time.NewTicker(frameDuration)
	canvas := newCanvas(config.CellsX, config.CellsY, config.WindowX, config.WindowY)

	started := false

	for range fpsClock.C {
		if win.Closed() {
			return
		}

		if !started {
			started = preStart(win, canvas)
		} else {
			canvas.Cells = config.Automaton.Step(canvas.Cells)
		}

		renderFrame(win, canvas, config.Automaton.GetColouring())
	}
}

func preStart(win *opengl.Window, canvas canvas) bool {
	if win.JustPressed(pixel.KeyS) {
		return true
	}

	if win.JustPressed(pixel.MouseButton1) {
		location := getVirtualPixelXY(win.MousePosition(), canvas)
		oldCell := canvas.Cells[uint(location.X)][uint(location.Y)]
		newCell := (oldCell + 1) % config.Automaton.CountStates()
		canvas.Cells[uint(location.X)][uint(location.Y)] = newCell
	}

	return false
}

func renderFrame(win *opengl.Window, canvas canvas, colourings []models.Rgb) {
	win.Canvas().SetPixels(canvas.paint(colourings))
	win.Update()
}

func premultiply(pixels []float64) []uint8 {
	premultiplied := make([]uint8, len(pixels))

	for i := 0; i < len(pixels); i += 4 {
		for j := 0; j < 3; j++ {
			premultiplied[i+j] = uint8(pixels[i+j] * 255)
		}
		premultiplied[i+3] = uint8(255)
	}

	return premultiplied
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

func (c canvas) paint(colourings []models.Rgb) []uint8 {
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

	return premultiply(pixels)
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

func getVirtualPixelXY(xy pixel.Vec, c canvas) pixel.Vec {
	x := uint(xy.X) / c.CellWidth
	y := uint(xy.Y) / c.CellHeight
	return pixel.Vec{X: float64(x), Y: float64(y)}
}

func getRealPixelIndex(realX, realY, canvasWidth uint) uint {
	return 4 * (realY*canvasWidth + realX)
}
