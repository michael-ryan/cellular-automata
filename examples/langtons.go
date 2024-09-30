package examples

import (
	"fmt"

	"github.com/michael-ryan/cellularautomata/v2/model"
)

// NewLangtons returns a [model.Automaton] that simulates Langton's ant.
// Ant colour and direction are encoded into eight separate states, all of which are painted red onscreen.
//
// When trying this automaton, remember to create at least one ant in edit mode first!
func NewLangtons() *model.Automaton {
	const (
		black = iota
		white
		blackAntN
		blackAntE
		blackAntS
		blackAntW
		whiteAntN
		whiteAntE
		whiteAntS
		whiteAntW
	)

	transitionSet := model.NewTransitionSet()

	// black transitions
	transitionSet.AddTransition(black, blackAntN, func(cell model.Cell) bool {
		// eastbound ant turning north
		neighbour, err := cell.Neighbour(-1, 0)
		if err != nil {
			return false
		}

		return neighbour == blackAntE || neighbour == whiteAntE
	})
	transitionSet.AddTransition(black, blackAntE, func(cell model.Cell) bool {
		// southbound ant turning east
		neighbour, err := cell.Neighbour(0, 1)
		if err != nil {
			return false
		}

		return neighbour == blackAntS || neighbour == whiteAntS
	})
	transitionSet.AddTransition(black, blackAntS, func(cell model.Cell) bool {
		// westbound ant turning south
		neighbour, err := cell.Neighbour(1, 0)
		if err != nil {
			return false
		}

		return neighbour == blackAntW || neighbour == whiteAntW
	})
	transitionSet.AddTransition(black, blackAntW, func(cell model.Cell) bool {
		// northbound ant turning west
		neighbour, err := cell.Neighbour(0, -1)
		if err != nil {
			return false
		}

		return neighbour == blackAntN || neighbour == whiteAntN
	})

	// white transitions
	transitionSet.AddTransition(white, whiteAntS, func(cell model.Cell) bool {
		// eastbound ant turning south
		neighbour, err := cell.Neighbour(-1, 0)
		if err != nil {
			return false
		}

		return neighbour == blackAntE || neighbour == whiteAntE
	})
	transitionSet.AddTransition(white, whiteAntW, func(cell model.Cell) bool {
		// southbound ant turning west
		neighbour, err := cell.Neighbour(0, 1)
		if err != nil {
			return false
		}

		return neighbour == blackAntS || neighbour == whiteAntS
	})
	transitionSet.AddTransition(white, whiteAntN, func(cell model.Cell) bool {
		// westbound ant turning north
		neighbour, err := cell.Neighbour(1, 0)
		if err != nil {
			return false
		}

		return neighbour == blackAntW || neighbour == whiteAntW
	})
	transitionSet.AddTransition(white, whiteAntE, func(cell model.Cell) bool {
		// northbound ant turning east
		neighbour, err := cell.Neighbour(0, -1)
		if err != nil {
			return false
		}

		return neighbour == blackAntN || neighbour == whiteAntN
	})

	for dir := blackAntN; dir <= blackAntW; dir++ {
		transitionSet.AddTransition(uint(dir), white, func(cell model.Cell) bool { return true })
	}

	for dir := whiteAntN; dir <= whiteAntW; dir++ {
		transitionSet.AddTransition(uint(dir), black, func(cell model.Cell) bool { return true })
	}

	colouring := make([]model.Rgb, 10)
	colouring[black] = model.Rgb{R: 0, G: 0, B: 0}
	colouring[white] = model.Rgb{R: 1, G: 1, B: 1}
	for state := blackAntN; state <= whiteAntW; state++ {
		colouring[state] = model.Rgb{R: 1, G: 0, B: 0}
	}

	automaton, err := model.NewAutomaton(transitionSet, colouring)
	if err != nil {
		panic(fmt.Errorf("unrecoverable error, something went wrong constructing Conway's Game of Life automaton: %w", err))
	}

	return automaton
}
