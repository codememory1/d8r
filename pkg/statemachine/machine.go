package statemachine

// Machine defines a state machine with named transitions between states.
type Machine[N comparable, S comparable] struct {
	transitions map[N]Transition[N, S]
}

// NewMachine creates a state machine from the provided transitions.
func NewMachine[N comparable, S comparable](transitions ...Transition[N, S]) *Machine[N, S] {
	items := make(map[N]Transition[N, S], len(transitions))

	for _, transition := range transitions {
		items[transition.Name] = transition
	}

	return &Machine[N, S]{
		transitions: items,
	}
}

// Can reports whether the named transition can be applied from the current state.
func (m *Machine[N, S]) Can(name N, current S) bool {
	transition, ok := m.transitions[name]

	if !ok {
		return false
	}

	return transition.Can(current)
}

// Transition resolves the target state for the named transition from the current state.
func (m *Machine[N, S]) Transition(name N, current S) (S, error) {
	transition, ok := m.transitions[name]

	if !ok {
		var zero S

		return zero, ErrTransitionNotFound
	}

	if !transition.Can(current) {
		var zero S

		return zero, InvalidTransitionError[N, S]{
			Name: name,
			From: current,
		}
	}

	return transition.To, nil
}
