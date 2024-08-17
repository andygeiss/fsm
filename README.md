<p align="center">
<img src="https://github.com/andygeiss/fsm/blob/main/logo.png?raw=true" />
</p>

# FSM - Finite State Machine

[![License](https://img.shields.io/github/license/andygeiss/fsm)](https://github.com/andygeiss/fsm/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/fsm)](https://github.com/andygeiss/fsm/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/fsm)](https://goreportcard.com/report/github.com/andygeiss/fsm)
[![Codacy Grade Badge](https://app.codacy.com/project/badge/Grade/57bb148a04154ae8b7ce40cecb78947c)](https://app.codacy.com/gh/andygeiss/fsm/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Cover Badge](https://app.codacy.com/project/badge/Coverage/57bb148a04154ae8b7ce40cecb78947c)](https://app.codacy.com/gh/andygeiss/fsm/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

Compute the next state with recursive state functions in Golang using generics and iterators.

## About

Based on Rob Pike's talk on [lexical scanning](https://www.youtube.com/watch?v=HxaD_trXwRE)
I thought about a version of a finite state machine (FSM) that uses
`Generics` from [go1.18](https://go.dev/blog/go1.18) and
Iterators” from [go1.23](https://go.dev/blog/go1.23), in order to
Golang's fantastic “batteries included” capabilities.

## Walkthrough

The best way to demonstrate the use of an FSM is to implement a game like “Super Mario”.
In this game, Mario changes his state and behavior depending on certain events,
as shown in the following illustration from the [Mario Wiki](https://www.mariowiki.com/Super_Mario_World):

<p align="center">
<img src="https://github.com/andygeiss/fsm/blob/main/mario.png?raw=true" />
</p>

Based on the image above, we could specify the `States` and `Events` as follows:

States:
1. Small Mario
2. Super Mario
3. Fire Mario
4. Cape Mario

Events:
1. Got Mushroom
2. Got Fire Flower
3. Got Feather

In Object-Oriented Programming (OOP), we would specify Mario
as an object that manages its internal/private state.
The behavior of Mario changes depending on the state
and is implemented as methods.
The game world knows its entities and must emit the events
based on player inputs, for example.

In Golang, however, we could implement each state as a function
that operates on data and returns a function (recursively).

### Initialize the world

Let's implement the first state transition from `Uninitialized`
to `Small Mario`.

```go
const (
    StateMarioUninitialized = iota
    StateMarioSmall
)

type World struct {
    eventCh    chan int
    marioState int
}

func smallMario(world *World) fsm.StateFn[*World] {
    world.marioState = StateMarioSmall
    return nil
}
```

Our implementation uses a channel to receive events from the game world.
Now we can execute the state machine to get the result as follows:

```go

// Create the inial world state.
world := &World{
    eventCh:    make(chan int, 0),
    marioState: StateMarioUninitialized,
}

// Start the state machine and wait for events.
doneCh := make(chan bool)
go fsm.Run(smallMario, world, doneCh)
<-doneCh

// State should be the value of StateMarioSmall = 1.
log.Printf("World marioState: %d", world.marioState)
```

### Mario got mushroom

After that, we will make the transition from `MarioSmall`
to `MarioSuper` after collecting a mushroom and receiving
the event `EventGotMushroom`.

```go
const (
    StateMarioUninitialized = iota
    StateMarioSmall
    StateMarioSuper // <-- NEW -->
)

func smallMario(world *World) fsm.StateFn[*World] {
    world.marioState = StateMarioSmall
    // <-- NEW
    switch <-world.eventCh {
    case EventGotMushroom:
        return superMario(world)
    }
    // -->
    return nil
}

// <-- NEW
func superMario(world *World) fsm.StateFn[*World] {
    world.marioState = StateMarioSuper
    return nil
}
// -->

// Start the state machine and wait for events.
doneCh := make(chan bool)
go fsm.Run(smallMario, world, doneCh)

// Emit events
world.eventCh <- EventGotMushroom // <-- NEW -->

// Wait for the end.
<-doneCh

// State should be the value of StateMarioSuper = 2.
log.Printf("World marioState: %d", world.marioState)
```
