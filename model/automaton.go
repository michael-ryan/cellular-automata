package model

import (
	"fmt"
	"sync"
)

// Rgb describes a single pixel's RGB colour value. All member fields must be in the closed interval [0-1].
type Rgb struct {
	R, G, B float64
}

// Automaton contains all the information needed to describe a cellular automaton. You should use the [NewAutomaton] function to create one.
type Automaton struct {
	colouring     []Rgb
	transitionSet TransitionSet
	states        uint
}

// at safely indexes the cell matrix c, returning an error if an off-grid value has been indexed.
func at(c [][]uint, x, y int) (uint, error) {
	if x < 0 || x >= len(c) || y < 0 || y >= len(c[0]) {
		return 0, fmt.Errorf("indexed off grid")
	}

	return c[x][y], nil
}

// NewAutomaton constructs a new Automaton from a given TransitionSet and colouring setup.
// For a state n, transitions[n] should describe the transition rules and colouring[n] should define its render colour.
// It is ill-advised to set up the [TransitionSet] yourself. Rather, you should use [NewTransitionSet] and [TransitionSet.AddTransition].
//
// Examples of how to set up an automaton are available in the [github.com/michael-ryan/cellularautomata/examples] package.
func NewAutomaton(transitions TransitionSet, colouring []Rgb) (*Automaton, error) {
	if len(transitions) != len(colouring) {
		return nil, fmt.Errorf("mismatched lengths of transitions and colouring: %v != %v", len(transitions), len(colouring))
	}

	if len(transitions) <= 1 {
		return nil, fmt.Errorf("it does not make sense to create an automaton that describes 0 or 1 states")
	}

	bad := false
	for state, rgb := range colouring {
		bad = bad || rgb.R > 1 || rgb.R < 0
		bad = bad || rgb.G > 1 || rgb.G < 0
		bad = bad || rgb.B > 1 || rgb.B < 0

		if bad {
			return nil, fmt.Errorf("colouring rule at index %v invalid, all values must be in the closed interval [0-1]: %+v", state, rgb)
		}
	}

	states := uint(len(transitions))

	for fromState, stateTransitions := range transitions {
		for ruleNumber, t := range stateTransitions {
			if t.NewState >= states {
				return nil, fmt.Errorf("state index %v, rule index %v has invalid new state %v (max = len(transitions) - 1 = %v)", fromState, ruleNumber, t.NewState, len(transitions)-1)
			}
		}
	}

	return &Automaton{states: states,
		transitionSet: transitions,
		colouring:     colouring,
	}, nil
}

// States describes the number of states defined in this automaton.
func (a Automaton) CountStates() uint {
	return a.states
}

// Colouring is an array of RGB values, instructing the renderer what colour to show a given state. State n should have its colour described in Colouring[n].
func (a Automaton) GetColouring() []Rgb {
	colouringCopy := make([]Rgb, len(a.colouring))
	copy(colouringCopy, a.colouring)
	return colouringCopy
}

// Step simulates a single time step.
// All cells will have their transition rules checked, and a new array is returned representing the new states of all the cells.
//
// This is a pure function, and will not modify any state, so it is safe to call manually.
// However, it is not needed to call this manually if using the provided graphical rendering package [github.com/michael-ryan/cellularautomata/].
func (a Automaton) Step(c [][]uint) [][]uint {
	new := make([][]uint, len(c))
	for x := range len(new) {
		new[x] = make([]uint, len(c[0]))
		for y := range len(new[x]) {
			new[x][y] = c[x][y]
		}
	}

	type edit struct {
		x, y     int
		newState uint
	}

	wg := sync.WaitGroup{}
	editChan := make(chan edit)

	for x := range len(c) {
		for y := range len(c[0]) {
			// run each cell compute in its own goroutine
			wg.Add(1)
			go func(x, y int, c [][]uint, model TransitionSet, editChan chan<- edit) {
				defer wg.Done()
				thisCell, err := at(c, x, y)
				if err != nil {
					panic("Something has gone very wrong. We indexed outside of the grid in the step function.")
				}

				cell := Cell{
					x:     x,
					y:     y,
					cells: c,
				}

				for _, t := range model[thisCell] {
					if t.Predicate(cell) {
						editChan <- edit{
							x:        x,
							y:        y,
							newState: t.NewState,
						}
						return
					}
				}
			}(x, y, c, a.transitionSet, editChan)
		}
	}

	go func() {
		wg.Wait()
		close(editChan)
	}()

	for edit := range editChan {
		new[edit.x][edit.y] = edit.newState
	}

	return new
}
