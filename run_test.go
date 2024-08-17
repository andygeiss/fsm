package fsm_test

import (
	"testing"

	"github.com/andygeiss/fsm"
)

type Cfg struct {
	Value int
}

func TestRun_Transition_To_New_State_Should_Set_Value_To_1(t *testing.T) {
	cfg := Cfg{Value: 0}
	stateFn := func(cfg *Cfg) fsm.StateFn[*Cfg] {
		cfg.Value = 1
		return nil
	}
	fsm.Run(stateFn, &cfg)
	if cfg.Value != 1 {
		t.Errorf("Value should be 1, but got %d", cfg.Value)
	}
}

func TestRun_Transition_From_State_A_To_State_B_Should_Increase_Value_To_2(t *testing.T) {
	cfg := Cfg{Value: 0}
	stateB := func(cfg *Cfg) fsm.StateFn[*Cfg] {
		cfg.Value++
		return nil
	}
	stateA := func(cfg *Cfg) fsm.StateFn[*Cfg] {
		cfg.Value = 1
		return stateB
	}
	fsm.Run(stateA, &cfg)
	if cfg.Value != 2 {
		t.Errorf("Value should be 2, but got %d", cfg.Value)
	}
}
