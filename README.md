# Cellular Automata Simulator

## Usage
Set up a `Config` object and supply it to `Launch`:
```Go
import (
    "github.com/michael-ryan/cellularAutomata/"
    "github.com/michael-ryan/cellularAutomata/models"
)

c := cellularautomata.Config{
    Fps:       15,
    CellsX:    72,
    CellsY:    128,
    WindowX:   720,
    WindowY:   1280,
    Automaton: models.NewForest(),
}

cellularautomata.Launch(c)
```

Note: Launch must be called from the main goroutine, due to a limitation in OpenGL.

You can easily construct your own automata. For an example, see the implementation for Conway's Game of Life [here](models/conways.go).

# Building

See requirements [here](https://github.com/gopxl/pixel?tab=readme-ov-file#requirements).