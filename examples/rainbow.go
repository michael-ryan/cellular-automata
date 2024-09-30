package examples

import (
	"fmt"

	"github.com/michael-ryan/cellularautomata/v2/model"
)

// TODO
func NewRainbow() *model.Automaton {
	const (
		black = iota
		red
		yellow
		green
		cyan
		blue
		purple
	)

	transitionSet := model.NewTransitionSet()

	transitionSet.AddTransition(black, red, func(cell model.Cell) bool {
		return cell.CountNeighbours(red, false) >= 1
	})

	truePredicate := func(cell model.Cell) bool { return true }

	transitionSet.AddTransition(red, yellow, truePredicate)
	transitionSet.AddTransition(yellow, green, truePredicate)
	transitionSet.AddTransition(green, cyan, truePredicate)
	transitionSet.AddTransition(cyan, blue, truePredicate)
	transitionSet.AddTransition(blue, purple, truePredicate)
	transitionSet.AddTransition(purple, red, truePredicate)

	colouring := make([]model.Rgb, 7)
	colouring[black] = model.Rgb{R: 0, G: 0, B: 0}
	colouring[red] = model.Rgb{R: 1, G: 0, B: 0}
	colouring[yellow] = model.Rgb{R: 1, G: 1, B: 0}
	colouring[green] = model.Rgb{R: 0, G: 1, B: 0}
	colouring[cyan] = model.Rgb{R: 0, G: 1, B: 1}
	colouring[blue] = model.Rgb{R: 0, G: 0, B: 1}
	colouring[purple] = model.Rgb{R: 1, G: 0, B: 1}

	automaton, err := model.NewAutomaton(transitionSet, colouring)
	if err != nil {
		panic(fmt.Errorf("unrecoverable error, something went wrong constructing Rainbow automaton: %w", err))
	}

	return automaton
}
