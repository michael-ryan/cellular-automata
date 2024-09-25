package models

import "fmt"

func NewConways() *Automaton {
	countAliveNeighbours := func(c [][]uint, x, y int) uint {
		neighbours := uint(0)
		for offsetX := -1; offsetX <= 1; offsetX++ {
			for offsetY := -1; offsetY <= 1; offsetY++ {
				if offsetY == 0 && offsetX == 0 {
					continue
				}

				neighbour, err := At(c, x+offsetX, y+offsetY)
				if err != nil {
					// gone off the grid, there's no neighbour here
					continue
				}

				if neighbour == 1 {
					neighbours++
				}
			}
		}
		return neighbours
	}

	const (
		dead = iota
		alive
	)

	transitionSet := NewTransitionSet()

	// dead rules
	transitionSet.AddTransition(dead, alive, func(c [][]uint, x, y int) bool {
		// 2 or 3 alive neighbours, become alive
		return countAliveNeighbours(c, x, y) == 3
	})

	// alive rules
	transitionSet.AddTransition(alive, dead, func(c [][]uint, x, y int) bool {
		// <2 alive neighbours, die
		return countAliveNeighbours(c, x, y) < 2
	})
	transitionSet.AddTransition(alive, dead, func(c [][]uint, x, y int) bool {
		// >3 alive neighbours, die
		return countAliveNeighbours(c, x, y) > 3
	})

	colouring := make([]Rgb, 2)
	colouring[dead] = Rgb{0, 0, 0}
	colouring[alive] = Rgb{1, 1, 1}

	automaton, err := NewAutomaton(transitionSet, colouring)
	if err != nil {
		panic(fmt.Errorf("unrecoverable error, something went wrong constructing Conway's Game of Life automaton: %w", err))
	}

	return automaton
}
