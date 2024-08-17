package fsm

import "context"

// StateFn is a recursive function which handles a specific state.
type StateFn[T any] func(ctx context.Context, cfg T) (stateFn StateFn[T])

// Run just run until the state goes nil.
func Run[T any](startFn StateFn[T], ctx context.Context, cfg T) {
	for state := startFn; state != nil; {
		state = state(ctx, cfg)
	}
}
