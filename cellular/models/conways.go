package models

type conways struct {
	model     [][]Transition
	colouring cellColouring
}

func (c *conways) GetTransitionModel() [][]Transition {
	return c.model
}

func (c *conways) GetCellColouring() cellColouring {
	return c.colouring
}

func (c *conways) Step() {
}

func NewConways() CellularAutomata {
	conways := conways{}

	model := make([][]Transition, 0)

	countAliveNeighbours := func(c Canvas, x, y int) uint {
		neighbours := uint(0)
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

				if neighbour == 1 {
					neighbours++
				}
			}
		}
		return neighbours
	}

	// dead
	state0 := make([]Transition, 0)
	state0 = append(state0, Transition{
		// 2 or 3 alive neighbours, become alive
		Predicate: func(c Canvas, x, y int) bool {
			return countAliveNeighbours(c, x, y) == 3
		},
		NewState: 1,
	})
	model = append(model, state0)

	state1 := make([]Transition, 0)
	state1 = append(state1, Transition{
		// <2 alive neighbours, die
		Predicate: func(c Canvas, x, y int) bool {
			return countAliveNeighbours(c, x, y) < 2
		},
		NewState: 0,
	})
	state1 = append(state1, Transition{
		// >3 alive neighbours, die
		Predicate: func(c Canvas, x, y int) bool {
			return countAliveNeighbours(c, x, y) > 3
		},
		NewState: 0,
	})
	state1 = append(state1, Transition{
		// 2 or 3 alive neighbours, live
		Predicate: func(c Canvas, x, y int) bool {
			alive := countAliveNeighbours(c, x, y)
			return alive == 2 || alive == 3
		},
		NewState: 1,
	})
	model = append(model, state1)

	conways.model = model
	conways.colouring = make(cellColouring, 2)
	conways.colouring[0] = Rgb{0, 0, 0}
	conways.colouring[1] = Rgb{1, 1, 1}

	return &conways
}
