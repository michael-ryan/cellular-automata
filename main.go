package main

import (
	"github.com/michael-ryan/cellular-automata/render"
	"github.com/michael-ryan/cellular-automata/render/models"
)

func main() {
	c := render.Config{
		Fps:       15,
		CellsX:    400,
		CellsY:    400,
		WindowX:   800,
		WindowY:   800,
		Automaton: models.NewForest(),
	}
	render.Launch(c)
}
