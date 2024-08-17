package fsm_test

import (
	"testing"

	"github.com/andygeiss/fsm"
)

const (
	EventGotMushroom = iota
)

const (
	StateMarioSmall = iota
	StateMarioSuper
)

type World struct {
	eventCh    chan int
	marioState int
}

func smallMario(world *World) fsm.StateFn[*World] {
	switch <-world.eventCh {
	case EventGotMushroom:
		return superMario(world)
	}
	return nil
}

func superMario(world *World) fsm.StateFn[*World] {
	world.marioState = StateMarioSuper
	return nil
}

func TestMario_SmallMario_To_SuperMario(t *testing.T) {
	// Arrange
	world := &World{
		eventCh:    make(chan int, 0),
		marioState: StateMarioSmall,
	}
	// Act
	go fsm.Run(smallMario, world)
	world.eventCh <- EventGotMushroom
	// Assert
	if world.marioState != StateMarioSuper {
		t.Errorf("Mario state should be %d, but got %d", StateMarioSuper, world.marioState)
	}
}
