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
`Generics` from [go1.18](https://go.dev/blog/go1.18) in order to
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

### Events, States and the World

Let's implement the events and states as follows:

```go
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
```

Our implementation will use a channel to receive events from the game world.
The world needs to know the state of Mario.

```go
type config struct {
    event      chan int
    stateMario int
}
```

### Mario got some stuff

We will implement the state transitions by using a `fsm.StateFn` as follows:

```go
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
```

This will result in a very clean approach that is easy to maintain (and test).
A complete example for Mario can be found [here](mario_test.go).
