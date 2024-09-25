package models

type langtons struct {
	model                    [][]Transition
	colouring                cellColouring
	direction, nextDirection uint8 // 0 is north, 1 east etc.

}

func (c *langtons) GetTransitionModel() [][]Transition {
	return c.model
}

func (c *langtons) GetCellColouring() cellColouring {
	return c.colouring
}

func (c *langtons) Step() {
	c.direction = c.nextDirection
}

func NewLangtons() CellularAutomata {
	langtons := langtons{direction: 0}

	black := Cell(0)
	white := Cell(1)
	blackAnt := Cell(2)
	whiteAnt := Cell(3)

	model := make([][]Transition, 0)

	blackTransitions := make([]Transition, 0)
	blackTransitions = append(blackTransitions, Transition{
		// At a black square, turn 90° counter-clockwise, flip the color of the square, move forward one unit
		Predicate: func(c Canvas, x, y int) bool {
			targetX := x
			targetY := y
			if langtons.direction == 0 {
				targetY--
			} else if langtons.direction == 1 {
				targetX--
			} else if langtons.direction == 2 {
				targetY++
			} else if langtons.direction == 3 {
				targetX++
			} else {
				panic("fuck you there is no direction")
			}

			cell, err := c.At(targetX, targetY)
			if err != nil {
				// out of bounds, skip
				return false
			}

			if cell == blackAnt || cell == whiteAnt {
				// dont overflow
				if langtons.direction == 0 {
					langtons.nextDirection = 3
				} else {
					langtons.nextDirection = langtons.direction - 1
				}
				return true
			}
			return false
		},
		NewState: blackAnt,
	})
	model = append(model, blackTransitions)

	whiteTransitions := make([]Transition, 0)
	whiteTransitions = append(whiteTransitions, Transition{
		// At a white square, turn 90° clockwise, flip the color of the square, move forward one unit
		Predicate: func(c Canvas, x, y int) bool {
			targetX := x
			targetY := y
			if langtons.direction == 0 {
				targetY--
			} else if langtons.direction == 1 {
				targetX--
			} else if langtons.direction == 2 {
				targetY++
			} else if langtons.direction == 3 {
				targetX++
			} else {
				panic("fuck you there is no direction")
			}

			cell, err := c.At(targetX, targetY)
			if err != nil {
				// out of bounds, skip
				return false
			}

			if cell == blackAnt || cell == whiteAnt {
				// dont overflow
				langtons.nextDirection = (langtons.direction + 1) % 4
				return true
			}
			return false
		},
		NewState: whiteAnt,
	})
	model = append(model, whiteTransitions)

	blackAntTransitions := make([]Transition, 0)
	blackAntTransitions = append(blackAntTransitions, Transition{
		Predicate: func(c Canvas, x, y int) bool {
			return true
		},
		NewState: white,
	})
	model = append(model, blackAntTransitions)

	whiteAntTransitions := make([]Transition, 0)
	whiteAntTransitions = append(whiteAntTransitions, Transition{
		Predicate: func(c Canvas, x, y int) bool {
			return true
		},
		NewState: black,
	})
	model = append(model, whiteAntTransitions)

	langtons.model = model

	langtons.colouring = make(cellColouring, 4)
	langtons.colouring[black] = Rgb{0, 0, 0}
	langtons.colouring[white] = Rgb{1, 1, 1}
	langtons.colouring[blackAnt] = Rgb{1, 0.1, 0.1}
	langtons.colouring[whiteAnt] = Rgb{0, 1, 0.4}

	return &langtons
}
