package model

import (
	"fmt"
	"math"
)

// Cell is provided as a parameter to the [Predicate] required for [TransitionSet.AddTransition].
// It exposes no fields, but it provides methods that allow the user to query the state of neighbouring cells.
type Cell struct {
	x, y  int
	cells [][]uint
}

// Neighbour checks the state of the neighbouring cell with a provided displacement.
// This function will return an error if a displacement of (0, 0) has been supplied, or there is no cell at that position (i.e. off the edge of the grid).
//
// Positive X goes right, positive Y goes up.
//
//	n, err = c.Neighbour(-1, 0) // neighbour to the left of c
//	n, err = c.Neighbour(1, -1) // neighbour to the bottom right of c
func (c Cell) Neighbour(x, y int) (uint, error) {
	if x == 0 && y == 0 {
		return 0, fmt.Errorf("cannot query self as neighbour")
	}

	return at(c.cells, c.x+x, c.y+y)
}

// CountNeighbours computes the number of neighbouring cells that have a given target state.
//
// Neighbour is defined by the moore parameter.
// If true, all eight surrounding cells are considered.
// If false, only the four orthogonally adjacent cells are considered.
//
// If this cell is at the edge of the grid, it will have fewer neighbours.
// Off-grid locations are considered to have no state, and will not contribute to the returned value.
func (c Cell) CountNeighbours(target uint, moore bool) uint {
	count := uint(0)
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			diagonal := math.Abs(float64(x))+math.Abs(float64(y)) == 2
			if !moore && diagonal {
				continue
			}

			neighbour, err := c.Neighbour(x, y)
			if err != nil {
				continue
			}

			if neighbour == target {
				count++
			}
		}
	}
	return count
}
