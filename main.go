// nolint
package main

import (
	pix "github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/michael-ryan/cellular-automata/cellular"
)

const (
	WIDTH  = 1280
	HEIGHT = 720
)

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Cellular Automata",
		Bounds: pix.R(0, 0, float64(WIDTH), float64(HEIGHT)),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	done := make(chan any)
	toDraw := make(chan []float64)

	cellular.Render(toDraw, done, 10, 128, 72, WIDTH, HEIGHT) // todo pull fps from the user level maybe

	for !win.Closed() {
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
