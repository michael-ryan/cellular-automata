package models

import (
	"fmt"
	"math/rand"
)

func NewForest() *Automaton {
	countNeighbours := func(c [][]uint, x, y int, target uint) uint {
		var count uint = 0
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

				if neighbour == target {
					count++
				}
			}
		}

		return count
	}

	const (
		dead = iota
		alive
		onFire
	)

	transitionSet := NewTransitionSet()

	// dead rules
	transitionSet.AddTransition(dead, alive, func(c [][]uint, y, x int) bool {
		// grow a random tree
		return rand.Float64() > 0.99999
	})
	transitionSet.AddTransition(dead, alive, func(c [][]uint, y, x int) bool {
		// grow a tree from a neighbouring tree
		for range countNeighbours(c, x, y, alive) {
			if rand.Float64() > 0.99 {
				return true
			}
		}

		return false
	})

	// alive rules
	transitionSet.AddTransition(alive, onFire, func(c [][]uint, y, x int) bool {
		// lightning sets a tree on fire
		return rand.Float64() > 0.9999
	})
	transitionSet.AddTransition(alive, onFire, func(c [][]uint, y, x int) bool {
		// neighbouring tree on fire, catch fire
		return countNeighbours(c, x, y, onFire) > 0 && rand.Float64() > 0.25
	})

	// burning rules
	transitionSet.AddTransition(onFire, dead, func(c [][]uint, y, x int) bool {
		// burn out
		return rand.Float64() > 0.3
	})

	colouring := make([]Rgb, 3)
	colouring[dead] = Rgb{0, 0, 0}
	colouring[alive] = Rgb{0, 1, 0}
	colouring[onFire] = Rgb{1, 0, 0}

	automaton, err := NewAutomaton(transitionSet, colouring)
	if err != nil {
		panic(fmt.Errorf("unrecoverable error, something went wrong constructing Conway's Game of Life automaton: %w", err))
	}

	return automaton
}
