package model

type Predicate func(cell Cell) bool

type transition struct {
	Predicate Predicate
	NewState  uint
}

type TransitionSet [][]transition

// NewTransitionSet creates a new [TransitionSet] for use with [NewAutomaton]
func NewTransitionSet() TransitionSet {
	return make(TransitionSet, 0)
}

// AddTransition adds a transition to this [TransitionSet]. States must be 0-indexed.
// It is recommended to set up your states as a const block, so they can be referred to by name:
//
//	const (
//		alive = iota
//		dying
//		dead
//	)
//
// Rule is a [Predicate] that is called with the cell in question.
// The [Cell] provides methods that can be used to decide on this cell's next state.
//
// The function should report whether the conditions are met for a state transition.
//
// For example, a transition rule where this cell turns from state 0 to state 1 if the cell to the right is in state 1 would look like this:
//
//	t.AddTransition(0, 1, func(cell Cell) bool {
//		right, err := cell.Neighbour(1, 0)
//		if err != nil {
//			return false
//		}
//		return right == 1
//	})
func (t *TransitionSet) AddTransition(fromState, toState uint, rule Predicate) {
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
