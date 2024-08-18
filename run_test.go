package fsm_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/fsm"
)

func setup() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*1)
}

func TestRun_Given_0_Should_Return_0_After_StateFn(t *testing.T) {
	// Arrange
	type stateData struct{ state int }
	ctx, cancel := setup()
	defer cancel()
	data := stateData{state: 0}
	stateFn := func(ctx context.Context, data stateData) fsm.StateFn[stateData] {
		return nil
	}

	// Act
	go fsm.Run(stateFn, ctx, data)

	// Assert
	select {
	case <-ctx.Done():
	}

	if data.state != 0 {
		t.Errorf("State should be 0, but got %d", data.state)
	}
}

func TestRun_Given_0_Should_Return_1_After_StateFn(t *testing.T) {
	// Arrange
	type stateData struct{ state int }
	ctx, cancel := setup()
	defer cancel()
	data := &stateData{state: 0}
	stateFn := func(ctx context.Context, data *stateData) fsm.StateFn[*stateData] {
		data.state = 1
		return nil
	}

	// Act
	go fsm.Run(stateFn, ctx, data)

	// Assert
	select {
	case <-ctx.Done():
	}

	if data.state != 1 {
		t.Errorf("State should be 1, but got %d", data.state)
	}
}

func TestRun_Given_0_Should_Return_2_After_Two_StateFn(t *testing.T) {
	// Arrange
	type stateData struct{ state int }
	ctx, cancel := setup()
	defer cancel()
	data := &stateData{state: 0}
	stateFnB := func(ctx context.Context, data *stateData) fsm.StateFn[*stateData] {
		data.state++
		return nil
	}
	stateFnA := func(ctx context.Context, data *stateData) fsm.StateFn[*stateData] {
		data.state++
		return stateFnB(ctx, data)
	}
	// Act
	go fsm.Run(stateFnA, ctx, data)

	// Assert
	select {
	case <-ctx.Done():
	}

	if data.state != 2 {
		t.Errorf("State should be 2, but got %d", data.state)
	}
}
