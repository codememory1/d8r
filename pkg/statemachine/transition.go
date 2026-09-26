package statemachine

// Transition defines a named transition from one or more source states
// to a target state.
type Transition[N comparable, S comparable] struct {
	Name N
	From []S
	To   S
}

// Can reports whether the transition can be applied from the given state.
func (t Transition[N, S]) Can(current S) bool {
	for _, from := range t.From {
		if from == current {
			return true
		}
	}

	return false
}

// T creates a named transition from the provided source states
// to the target state.
func T[N comparable, S comparable](name N, from []S, to S) Transition[N, S] {
	return Transition[N, S]{
		Name: name,
		From: from,
		To:   to,
	}
}
