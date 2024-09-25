# Cellular Automata Simulator

## Usage
Set up a `Config` object and supply it to `Launch`:
```Go
package main

import (
	"github.com/michael-ryan/cellularautomata"
	"github.com/michael-ryan/cellularautomata/models"
)

func main() {
	c := cellularautomata.Config{
		Fps:       15,
		CellsX:    128,
		CellsY:    72,
		WindowX:   1280,
		WindowY:   720,
		Automaton: models.NewForest(),
	}

	cellularautomata.Launch(c)
}
```

This will open a GUI window. You can click on cells to cycle their initial state, then press `S` on your keyboard to start the simulation.

Note: Launch must be called from the main goroutine, due to a limitation in OpenGL.

You can easily construct your own automata. For an example, see the implementation for Conway's Game of Life [here](models/conways.go).

## Building

See requirements [here](https://github.com/gopxl/pixel?tab=readme-ov-file#requirements).