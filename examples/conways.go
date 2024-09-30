package examples

import (
	"fmt"

	"github.com/michael-ryan/cellularautomata/model"
)

// NewConways returns a [model.Automaton] that represents Conway's Game of Life.
//
// [https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life]
func NewConways() *model.Automaton {
	const (
		dead = iota
		alive
	)

	transitionSet := model.NewTransitionSet()

	// dead rules
	transitionSet.AddTransition(dead, alive, func(c model.Cell) bool {
		// 3 alive neighbours, become alive
		return c.CountNeighbours(alive, true) == 3
	})

	// alive rules
	transitionSet.AddTransition(alive, dead, func(c model.Cell) bool {
		// <2 alive neighbours, die
		return c.CountNeighbours(alive, true) < 2
	})
	transitionSet.AddTransition(alive, dead, func(c model.Cell) bool {
		// >3 alive neighbours, die
		return c.CountNeighbours(alive, true) > 3
	})

	colouring := make([]model.Rgb, 2)
	colouring[dead] = model.Rgb{R: 0, G: 0, B: 0}
	colouring[alive] = model.Rgb{R: 1, G: 1, B: 1}

	automaton, err := model.NewAutomaton(transitionSet, colouring)
	if err != nil {
		panic(fmt.Errorf("unrecoverable error, something went wrong constructing Conway's Game of Life automaton: %w", err))
	}

	return automaton
}
