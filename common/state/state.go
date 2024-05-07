package state

// State Providers allow tree nodes to share data between each other by either writing
// the state value or reading from it.
type StateProvider[T any] struct {
	value T
}

func (p *StateProvider[T]) Get() T {
	return p.value
}

func (p *StateProvider[T]) Set(val T) {
	p.value = val
}

// Creates a state provider of a constant value that never changes
func MakeConstStateProvider[T any](val T) StateGetter[T] {
	return &constStateProvider[T]{value: val}
}

type constStateProvider[T any] struct {
	value T
}

func (p *constStateProvider[T]) Get() T {
	return p.value
}

type StateGetter[T any] interface {
	Get() T
}

type StateSetter[T any] interface {
	Set(val T)
}
