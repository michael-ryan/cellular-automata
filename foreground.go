package cellularautomata

import (
	"fmt"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/michael-ryan/cellularAutomata/models"
)

type Config struct {
	Fps, CellsX, CellsY, WindowX, WindowY uint
	Automaton                             *models.Automaton
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

	done := make(chan any)
	toDraw := make(chan []float64)
	clicks := make(chan coord, 1)
	start := make(chan any)

	startRenderer(toDraw, start, clicks, done, config)

	started := false

	for !win.Closed() {
		if win.JustPressed(pixel.MouseButton1) && !started {
			location := win.MousePosition()
			clicks <- [2]uint{uint(location.X), uint(location.Y)}
		} else if win.JustPressed(pixel.KeyS) && !started {
			close(clicks)
			start <- struct{}{}
			started = true
		}

		draw, ok := <-toDraw
		if !ok {
			return
		}
		win.Canvas().SetPixels(premultiply(draw))
		win.Update()
	}
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
