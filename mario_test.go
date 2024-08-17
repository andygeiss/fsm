package fsm_test

import (
	"testing"

	"github.com/andygeiss/fsm"
)

const (
	EventGoFeather = iota
	EventGoFireFlower
	EventGotMushroom
)

const (
	StateUndefined = iota
	StateMarioCape
	StateMarioFire
	StateMarioSmall
	StateMarioSuper
)

type World struct {
	eventCh    chan int
	marioState int
}

func fireMario(world *World) fsm.StateFn[*World] {
	world.marioState = StateMarioFire
	return nil
}

func smallMario(world *World) fsm.StateFn[*World] {
	world.marioState = StateMarioSmall
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
		marioState: StateUndefined,
	}
	// Act
	doneCh := make(chan bool)
	go fsm.Run(smallMario, world, doneCh)
	world.eventCh <- EventGotMushroom
	<-doneCh
	// Assert
	if world.marioState != StateMarioSuper {
		t.Errorf("Mario state should be %d, but got %d", StateMarioSuper, world.marioState)
	}
}
