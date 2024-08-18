package fsm

import "context"

// StateFn is a recursive function which handles a specific state.
type StateFn[T any] func(ctx context.Context, stateData T) (stateFn StateFn[T])

// Run just run until the state goes nil.
func Run[T any](startFn StateFn[T], ctx context.Context, stateData T) {
	for state := startFn; state != nil; {
		state = state(ctx, stateData)
	}
}
