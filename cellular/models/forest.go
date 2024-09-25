package models

import (
	"math/rand"
)

type forest struct {
	model     [][]Transition
	colouring cellColouring
}

func (c *forest) GetTransitionModel() [][]Transition {
	return c.model
}

func (c *forest) GetCellColouring() cellColouring {
	return c.colouring
}

func (c *forest) Step() {
}

func NewForest() CellularAutomata {
	forest := forest{}

	dead := Cell(0)
	alive := Cell(1)
	onFire := Cell(2)

	model := make([][]Transition, 0)

	countNeighbours := func(c Canvas, x, y int, target Cell) uint {
		var count uint = 0
		for offsetX := -1; offsetX <= 1; offsetX++ {
			for offsetY := -1; offsetY <= 1; offsetY++ {
				if offsetY == 0 && offsetX == 0 {
					continue
				}

				neighbour, err := c.At(x+offsetX, y+offsetY)
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

	// dead
	state0 := make([]Transition, 0)
	state0 = append(state0, Transition{
		// grow a random tree
		Predicate: func(c Canvas, x, y int) bool {
			return rand.Float64() > 0.99999
		},
		NewState: alive,
	})
	state0 = append(state0, Transition{
		// grow a tree from a neighbouring tree
		Predicate: func(c Canvas, x, y int) bool {
			for range countNeighbours(c, x, y, alive) {
				if rand.Float64() > 0.99 {
					return true
				}
			}

			return false
		},
		NewState: alive,
	})
	model = append(model, state0)

	// alive
	state1 := make([]Transition, 0)
	state1 = append(state1, Transition{
		// lightning sets a tree on fire
		Predicate: func(c Canvas, x, y int) bool {
			return rand.Float64() > 0.9999
		},
		NewState: onFire,
	})
	state1 = append(state1, Transition{
		// neighbouring tree on fire, catch fire
		Predicate: func(c Canvas, x, y int) bool {
			return countNeighbours(c, x, y, onFire) > 0 && rand.Float64() > 0.25
		},
		NewState: onFire,
	})
	model = append(model, state1)

	// burning
	state2 := make([]Transition, 0)
	state2 = append(state2, Transition{
		// burn out
		Predicate: func(c Canvas, x, y int) bool {
			return rand.Float64() > 0.3
		},
		NewState: dead,
	})
	model = append(model, state2)

	forest.model = model

	forest.colouring = make(cellColouring, 3)
	forest.colouring[dead] = Rgb{0, 0, 0}
	forest.colouring[alive] = Rgb{0, 1, 0}
	forest.colouring[onFire] = Rgb{1, 0, 0}

	return &forest
}
