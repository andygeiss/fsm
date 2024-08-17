package fsm

// StateFn is a recursive type which operates on data and returns a StateFn.
type StateFn[T any] func(data T) (stateFn StateFn[T])

// Run just run until the state goes nil.
func Run[T any](startFn StateFn[T], data T, doneCh chan bool) {
	for state := startFn; state != nil; {
		state = state(data)
	}
	doneCh <- true
}
