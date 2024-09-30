package examples

import (
	"fmt"
	"math/rand"

	"github.com/michael-ryan/cellularautomata/v2/model"
)

// NewForest returns a [model.Automaton] that simulates growing forests.
//
// Cells are either dead (black), alive (green) or on fire (red).
//   - If a dead cell has a neighbouring alive cell, it may turn into an alive cell. The more neighbouring alive cells, the higher the chance. There is also a very very low chance of a dead -> alive transition with no alive neighbours.
//   - An alive cell may randomly catch fire with a very low chance. It also has a high chance of spreading fire from a neighbouring cell.
//   - A cell that is on fire will (eventually) transition to the dead state, with a 70% chance rolled on each time step.
func NewForest() *model.Automaton {
	const (
		dead = iota
		alive
		onFire
	)

	transitionSet := model.NewTransitionSet()

	// dead rules
	transitionSet.AddTransition(dead, alive, func(cell model.Cell) bool {
		// grow a random tree
		return rand.Float64() > 0.99999
	})
	transitionSet.AddTransition(dead, alive, func(cell model.Cell) bool {
		// grow a tree from a neighbouring tree
		for range cell.CountNeighbours(alive, true) {
			if rand.Float64() > 0.99 {
				return true
			}
		}

		return false
	})

	// alive rules
	transitionSet.AddTransition(alive, onFire, func(cell model.Cell) bool {
		// lightning sets a tree on fire
		return rand.Float64() > 0.9999
	})
	transitionSet.AddTransition(alive, onFire, func(cell model.Cell) bool {
		// neighbouring tree on fire, catch fire
		return cell.CountNeighbours(onFire, true) > 0 && rand.Float64() > 0.25
	})

	// burning rules
	transitionSet.AddTransition(onFire, dead, func(cell model.Cell) bool {
		// burn out
		return rand.Float64() > 0.3
	})

	colouring := make([]model.Rgb, 3)
	colouring[dead] = model.Rgb{R: 0, G: 0, B: 0}
	colouring[alive] = model.Rgb{R: 0, G: 1, B: 0}
	colouring[onFire] = model.Rgb{R: 1, G: 0, B: 0}

	automaton, err := model.NewAutomaton(transitionSet, colouring)
	if err != nil {
		panic(fmt.Errorf("unrecoverable error, something went wrong constructing Conway's Game of Life automaton: %w", err))
	}

	return automaton
}
