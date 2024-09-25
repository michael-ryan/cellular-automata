package models

type predicate func(c [][]uint, x int, y int) bool

type transition struct {
	Predicate predicate
	NewState  uint
}

type transitionSet [][]transition

func NewTransitionSet() transitionSet {
	return make(transitionSet, 0)
}

// AddTransition adds a transition to this transitionSet. States are 0-indexed.
func (t *transitionSet) AddTransition(fromState, toState uint, rule predicate) {
	maxIndex := len(*t) - 1
	if maxIndex < int(fromState) {
		for range int(fromState) - maxIndex {
			*t = append(*t, make([]transition, 0))
		}
	}

	(*t)[fromState] = append((*t)[fromState], transition{
		NewState:  toState,
		Predicate: rule,
	})
}
