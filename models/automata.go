package models

import (
	"fmt"
	"sync"
)

// Rgb describes a single pixel's RGB colour value. All member fields must be in the closed interval [0-1].
type Rgb struct {
	R, G, B float64
}

// At safely indexes the cell matrix c, returning an error if an off-grid value has been indexed.
func At(c [][]uint, x, y int) (uint, error) {
	if x < 0 || x >= len(c) {
		return 0, fmt.Errorf("index x=%v off grid", x)
	}

	if y < 0 || y >= len(c[0]) {
		return 0, fmt.Errorf("index y=%v off grid", y)
	}

	return c[x][y], nil
}

// Automaton contains all the information needed to describe a cellular automaton. You should use the [NewAutomaton] function to create one.
type Automaton struct {
	colouring     []Rgb
	transitionSet transitionSet
	states        uint
}

func NewAutomaton(transitions transitionSet, colouring []Rgb) (*Automaton, error) {
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

// TransitionSet is an array of arrays of transition rules. All transition rules that state n can undergo should be in Transitions[n].
// It is acceptable for any inner array to be an empty array, to denote a dead-end state.
func (a Automaton) GetTransitionSet() transitionSet {
	transitionsCopy := make(transitionSet, len(a.transitionSet))
	for i := range a.transitionSet {
		transitionsCopy[i] = make([]transition, len(a.transitionSet[i]))
		copy(transitionsCopy[i], a.transitionSet[i])
	}
	return transitionsCopy
}

// Step simulates a single step. Typically this would not be called manually.
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
			go func(x, y int, c [][]uint, model transitionSet, editChan chan<- edit) {
				defer wg.Done()
				thisCell, err := At(c, x, y)
				if err != nil {
					panic("Something has gone very wrong. We indexed outside of the grid in the step function.")
				}
				for _, t := range model[thisCell] {
					if t.Predicate(c, x, y) {
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
