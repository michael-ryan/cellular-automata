package cellularautomata

import (
	"fmt"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/michael-ryan/cellularautomata/v2/model"
)

type Config struct {
	// Fps defines the target FPS of the simulation. This is the number of simulated time steps, per second.
	//
	// If Fps is sufficiently high, the simulation will simply run as fast as the hardware allows.
	Fps uint
	// CellsX and CellsY define the dimensions of the grid of cells.
	CellsX, CellsY uint
	// WindowX and WindowY define the number of pixels the GUI canvas should span.
	//
	// These values should equal or exceed CellsX and CellsY respectively.
	WindowX, WindowY uint
	// Automaton defines the cell states, their colours and their transition rules.
	//
	// For examples of how to define an Automaton, see the examples package: [github.com/michael-ryan/cellularautomata/examples]
	Automaton *model.Automaton
	// InitialState defines the initial state of all cells on the grid.
	InitialState uint
	// SkipEditor denotes whether to skip the initial edit mode of the grid. If false, the program will launch in edit mode, and the user can click on cells to cycle their initial state.
	// Pressing S on the keyboard will start the simulation.
	SkipEditor bool
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

var configChan chan Config = make(chan Config, 1)

// Launch opens a GUI window that renders the simulation.
// This function will block until the window is closed.
func Launch(config Config) error {
	if config.CellsX > config.WindowX {
		return fmt.Errorf("cellsX (%v) cannot be larger than windowX (%v), since each cell requires at least one pixel", config.CellsX, config.WindowX)
	}

	if config.CellsX > config.WindowX {
		return fmt.Errorf("cellsY (%v) cannot be larger than windowY (%v), since each cell requires at least one pixel", config.CellsY, config.WindowY)
	}

	stateCount := len(config.Automaton.GetColouring())
	if int(config.InitialState) >= stateCount {
		return fmt.Errorf("initialState too high at %v, there are only %v states defined, so initialState is bounded by [0-%v]", config.InitialState, stateCount, stateCount-1)
	}

	// dirty hack to get parameters into the opengl.Run callback
	configChan <- config
	opengl.Run(launch)

	return nil
}

func launch() {
	config := <-configChan

	windowConfig := opengl.WindowConfig{
		Title:  "Cellular Automata",
		Bounds: pixel.R(0, 0, float64(config.WindowX), float64(config.WindowY)),
		VSync:  true,
	}

	win, err := opengl.NewWindow(windowConfig)
	if err != nil {
		panic(err)
	}

	frameDuration := time.Second / time.Duration(config.Fps)
	fpsClock := time.NewTicker(frameDuration)

	canvas := newCanvas(config.CellsX, config.CellsY, config.WindowX, config.WindowY, config.InitialState)

	started := config.SkipEditor

	for range fpsClock.C {
		if win.Closed() {
			return
		}

		if !started {
			started = preStart(win, canvas, config.Automaton.CountStates())
		} else {
			canvas.Cells = config.Automaton.Step(canvas.Cells)
		}

		renderFrame(win, canvas, config.Automaton.GetColouring())
	}
}

func preStart(win *opengl.Window, canvas canvas, stateCount uint) bool {
	if win.JustPressed(pixel.KeyS) {
		return true
	}

	if win.JustPressed(pixel.MouseButton1) {
		location := getVirtualPixelXY(win.MousePosition(), canvas)
		oldCell := canvas.Cells[uint(location.X)][uint(location.Y)]
		newCell := (oldCell + 1) % stateCount
		canvas.Cells[uint(location.X)][uint(location.Y)] = newCell
	}

	return false
}

func renderFrame(win *opengl.Window, canvas canvas, colourings []model.Rgb) {
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

func newCanvas(width, height, realWidth, realHeight, initialState uint) canvas {
	c := canvas{}
	c.Cells = make([][]uint, width)
	for x := range c.Cells {
		c.Cells[x] = make([]uint, height)
		for y := range c.Cells[x] {
			c.Cells[x][y] = initialState
		}
	}
	c.Width = width
	c.Height = height
	c.RealWidth = realWidth
	c.RealHeight = realHeight

	c.CellWidth = c.RealWidth / c.Width
	c.CellHeight = c.RealHeight / c.Height

	return c
}

func (c canvas) paint(colourings []model.Rgb) []uint8 {
	pixels := make([]float64, 4*c.RealWidth*c.RealHeight)

	for x := range c.Width {
		for y := range c.Height {
			pixels = setPixel(pixels, uint(x), uint(y), colourings[c.Cells[x][y]], c)
		}
	}

	return premultiply(pixels)
}

func setPixel(pixels []float64, x, y uint, colour model.Rgb, c canvas) []float64 {
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
