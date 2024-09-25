package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/michael-ryan/cellular-automata/cellular"
)

const (
	WIDTH  = 1280
	HEIGHT = 720
)

func main() {
	opengl.Run(run)
}

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Cellular Automata",
		Bounds: pixel.R(0, 0, float64(WIDTH), float64(HEIGHT)),
		VSync:  true,
	}

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	done := make(chan any)
	toDraw := make(chan []float64)
	clicks := make(chan cellular.Coord, 1)
	start := make(chan any)

	cellular.StartRenderer(toDraw, start, clicks, done, 30, WIDTH/2, HEIGHT/2, WIDTH, HEIGHT) // todo pull fps from the user level maybe

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
