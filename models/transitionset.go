package models

type predicate func(c [][]uint, x int, y int) bool

type transition struct {
	Predicate predicate
	NewState  uint
}

type TransitionSet [][]transition

// NewTransitionSet creates a new [TransitionSet] for use with [NewAutomaton]
func NewTransitionSet() TransitionSet {
	return make(TransitionSet, 0)
}

// AddTransition adds a transition to this [TransitionSet]. States are 0-indexed.
// Rule is a func([][]uint, int int) bool that is called with:
//   - c [][]uint: A 2D array that describes the all cell states, indexable as c[x][y] where c[0][0] is the bottom left cell
//   - x: the cell in question's x coordinate
//   - y: the cell in question's y coordinate
//
// The function should report whether the conditions are met for a state transition.
//
// It is encouraged to use the [At] function to safely index the cell states.
//
// For example, a transition rule where this cell turns from state 0 to state 1 if the cell to the right is in state 1 would look like this:
//
//	t.AddTransition(0, 1, func(c [][]uint, x int, y int) bool {
//		target, err := At(c, x+1, y)
//		if err != nil {
//			return false
//		}
//		return target == 1
//	})
func (t *TransitionSet) AddTransition(fromState, toState uint, rule predicate) {
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
