package models

import "fmt"

func NewLangtons() *Automaton {
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

	transitionSet := NewTransitionSet()

	// black transitions
	transitionSet.AddTransition(black, blackAntN, func(c [][]uint, x, y int) bool {
		// eastbound ant turning north
		cell, err := At(c, x-1, y)
		if err != nil {
			return false
		}

		return cell == blackAntE || cell == whiteAntE
	})
	transitionSet.AddTransition(black, blackAntE, func(c [][]uint, x, y int) bool {
		// southbound ant turning east
		cell, err := At(c, x, y+1)
		if err != nil {
			return false
		}

		return cell == blackAntS || cell == whiteAntS
	})
	transitionSet.AddTransition(black, blackAntS, func(c [][]uint, x, y int) bool {
		// westbound ant turning south
		cell, err := At(c, x+1, y)
		if err != nil {
			return false
		}

		return cell == blackAntW || cell == whiteAntW
	})
	transitionSet.AddTransition(black, blackAntW, func(c [][]uint, x, y int) bool {
		// northbound ant turning west
		cell, err := At(c, x, y-1)
		if err != nil {
			return false
		}

		return cell == blackAntN || cell == whiteAntN
	})

	// white transitions
	transitionSet.AddTransition(white, whiteAntS, func(c [][]uint, x, y int) bool {
		// eastbound ant turning south
		cell, err := At(c, x-1, y)
		if err != nil {
			return false
		}

		return cell == blackAntE || cell == whiteAntE
	})
	transitionSet.AddTransition(white, whiteAntW, func(c [][]uint, x, y int) bool {
		// southbound ant turning west
		cell, err := At(c, x, y+1)
		if err != nil {
			return false
		}

		return cell == blackAntS || cell == whiteAntS
	})
	transitionSet.AddTransition(white, whiteAntN, func(c [][]uint, x, y int) bool {
		// westbound ant turning north
		cell, err := At(c, x+1, y)
		if err != nil {
			return false
		}

		return cell == blackAntW || cell == whiteAntW
	})
	transitionSet.AddTransition(white, whiteAntE, func(c [][]uint, x, y int) bool {
		// northbound ant turning east
		cell, err := At(c, x, y-1)
		if err != nil {
			return false
		}

		return cell == blackAntN || cell == whiteAntN
	})

	for dir := blackAntN; dir <= blackAntW; dir++ {
		transitionSet.AddTransition(uint(dir), white, func(c [][]uint, x, y int) bool { return true })
	}

	for dir := whiteAntN; dir <= whiteAntW; dir++ {
		transitionSet.AddTransition(uint(dir), black, func(c [][]uint, x, y int) bool { return true })
	}

	colouring := make([]Rgb, 10)
	colouring[black] = Rgb{0, 0, 0}
	colouring[white] = Rgb{1, 1, 1}
	for state := blackAntN; state <= whiteAntW; state++ {
		colouring[state] = Rgb{1, 0, 0}
	}

	automaton, err := NewAutomaton(transitionSet, colouring)
	if err != nil {
		panic(fmt.Errorf("unrecoverable error, something went wrong constructing Conway's Game of Life automaton: %w", err))
	}

	return automaton
}
