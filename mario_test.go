package fsm_test

import (
	"context"
	"testing"

	"github.com/andygeiss/fsm"
)

const (
	EventGotFeather = iota
	EventGotFireFlower
	EventGotMushroom
)

const (
	StateMarioCape = iota
	StateMarioFire
	StateMarioSmall
	StateMarioSuper
	StateMarioUndefined
)

type config struct {
	event      chan int
	stateMario int
}

func capeMario(ctx context.Context, cfg *config) fsm.StateFn[*config] {
	cfg.stateMario = StateMarioCape
	select {
	case event := <-cfg.event:
		switch event {
		case EventGotFireFlower:
			return fireMario(ctx, cfg)
		}
	case <-ctx.Done():
		return nil
	}
	return nil
}

func fireMario(ctx context.Context, cfg *config) fsm.StateFn[*config] {
	cfg.stateMario = StateMarioFire
	select {
	case event := <-cfg.event:
		switch event {
		case EventGotFeather:
			return capeMario(ctx, cfg)
		}
	case <-ctx.Done():
		return nil
	}
	return nil
}

func smallMario(ctx context.Context, cfg *config) fsm.StateFn[*config] {
	cfg.stateMario = StateMarioSmall
	select {
	case event := <-cfg.event:
		switch event {
		case EventGotFeather:
			return capeMario(ctx, cfg)
		case EventGotFireFlower:
			return fireMario(ctx, cfg)
		case EventGotMushroom:
			return superMario(ctx, cfg)
		}
	case <-ctx.Done():
		return nil
	}
	return nil
}

func superMario(ctx context.Context, cfg *config) fsm.StateFn[*config] {
	cfg.stateMario = StateMarioSuper
	select {
	case event := <-cfg.event:
		switch event {
		case EventGotFeather:
			return capeMario(ctx, cfg)
		case EventGotFireFlower:
			return fireMario(ctx, cfg)
		}
	case <-ctx.Done():
		return nil
	}
	return nil
}

func TestMario(t *testing.T) {
	testcases := []struct {
		desc        string
		startFn     fsm.StateFn[*config]
		startState  int
		event       int
		targetState int
	}{
		{desc: "Given smallMario When EventGotMushroom Should Be StateMarioSuper",
			startFn:     smallMario,
			startState:  StateMarioUndefined,
			event:       EventGotMushroom,
			targetState: StateMarioSuper,
		},
		{desc: "Given smallMario When EventGotFeather Should Be StateMarioCape",
			startFn:     smallMario,
			startState:  StateMarioUndefined,
			event:       EventGotFeather,
			targetState: StateMarioCape,
		},
		{desc: "Given smallMario When EventGotFireFlower Should Be StateMarioFire",
			startFn:     smallMario,
			startState:  StateMarioUndefined,
			event:       EventGotFireFlower,
			targetState: StateMarioFire,
		},
		{desc: "Given superMario When EventGotFeather Should Be StateMarioCape",
			startFn:     superMario,
			startState:  StateMarioUndefined,
			event:       EventGotFeather,
			targetState: StateMarioCape,
		},
		{desc: "Given superMario When EventGotFireFlower Should Be StateMarioFire",
			startFn:     superMario,
			startState:  StateMarioUndefined,
			event:       EventGotFireFlower,
			targetState: StateMarioFire,
		},
		{desc: "Given fireMario When EventGotFeather Should Be StateMarioCape",
			startFn:     fireMario,
			startState:  StateMarioUndefined,
			event:       EventGotFeather,
			targetState: StateMarioCape,
		},
		{desc: "Given capeMario When EventGotFireFlower Should Be StateMarioFire",
			startFn:     capeMario,
			startState:  StateMarioUndefined,
			event:       EventGotFireFlower,
			targetState: StateMarioFire,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.desc, func(t *testing.T) {

			ctx, cancel := setup()
			defer cancel()

			cfg := &config{event: make(chan int), stateMario: testcase.startState}

			go fsm.Run(testcase.startFn, ctx, cfg)
			go func() { cfg.event <- testcase.event }()

			select {
			case <-ctx.Done():
			}

			if cfg.stateMario != testcase.targetState {
				t.Errorf("stateMario should be %d, but got %d", testcase.targetState, cfg.stateMario)
			}
		})
	}
}
