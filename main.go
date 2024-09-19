package main

import (
	"math"
	"time"

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

	cellular.Render(toDraw, done, 10, 72, 128, HEIGHT, WIDTH) // todo pull fps from the user level maybe

	for !win.Closed() {
		select {
		case draw, ok := <-toDraw:
			if !ok {
				return
			}
			win.Canvas().SetPixels(premultiply(draw))
			win.Update()
		}
	}
}

func render(drawChan chan<- []float64, doneChan <-chan any, fps uint) {
	// barf shit down drawchan to render it

	// fps^-1 * 1000 = milliseconds per frame
	// 10fps -> (10)^-1 * 1000 = 0.1 * 1000 = 100
	frameDuration := time.Duration(math.Pow(float64(fps), -1)) * time.Millisecond

	fpsClock := time.NewTicker(frameDuration)

	for {
		select {
		case <-fpsClock.C:
			drawChan <- make([]float64, 4*WIDTH*HEIGHT)
		case <-doneChan:
			close(drawChan)
			return
		}
	}
}

func premultiply(pixels []float64) []uint8 {
	premultiplied := make([]uint8, len(pixels))

	for i := 0; i < len(pixels)/4; i++ {
		alpha := pixels[i+3]
		for j := 0; j < 3; j++ {
			premultiplied[i+j] = uint8(255 * (pixels[i+j] * alpha))
		}
		premultiplied[i+3] = uint8(255 * alpha)
	}

	return premultiplied
}
