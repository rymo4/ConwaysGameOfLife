package main

import (
	"github.com/rymo4/life/universe"
	"testing"
)

type roundtrip struct {
	Input    string
	Answer   string
	Toroidal bool
}

func (r *roundtrip) Check() (bool, string) {
	u := universe.LoadFromString(r.Input)
	u.Toroidal = r.Toroidal
	u.Next()
	return u.String() == r.Answer, u.String()
}

type testsuite []roundtrip

func (ts testsuite) Run(t *testing.T) {
	for _, test := range ts {
		ok, guess := test.Check()
		if !ok {
			t.Errorf("States do not match:\n%s\n%s\n Initial: %s", test.Answer, guess, test.Input)
		}
	}
}

func TestCanonicalStringRoundtrip(t *testing.T) {
	u := universe.LoadFromCanonicalString("3,2|0,0,")
	if u.Width != 3 || u.Height != 2 {
		t.Error("Size does not match")
	}

	tests := make(testsuite, 0)
	tests = append(tests, roundtrip{Input: "...\nOOO\n...\n", Answer: ".O.\n.O.\n.O.\n"})
	tests = append(tests, roundtrip{Input: ".O.\n.O.\n.O.\n", Answer: "...\nOOO\n...\n"})
	tests = append(tests, roundtrip{Input: "...\n...\n...\n", Answer: "...\n...\n...\n"})
	tests = append(tests, roundtrip{Input: "..\n..\n..\n", Answer: "..\n..\n..\n"})
	tests = append(tests, roundtrip{Input: ".OO.\n....\n.OO.\n", Answer: ".OO.\n....\n.OO.\n", Toroidal: true})
	tests.Run(t)
}
