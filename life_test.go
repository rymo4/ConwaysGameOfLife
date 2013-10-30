package main

import (
	"github.com/rymo4/life/universe"
	"testing"
)

// Possible data flows:
// 1. CanonicalString to Universe
// 2. Universe to CanonicalString
// 3. MapString to Universe
// 4. Universe to MapString
// 5. File to Universe
// Maybe later: 6. Universe to file
// Tests must send data through all pathes

//-------------------
// MapString -> Universe -> next Universe -> MapString
// Tests: 3, 4
//-------------------
type mapToNextMap struct {
	Input, Answer string
	Toroidal      bool
}

func (m mapToNextMap) Check() (bool, string) {
	u := universe.LoadFromString(m.Input)
	u.Toroidal = m.Toroidal
	u.Next()
	message := "Answer:\n" + m.Answer + "\nGuess:\n" + u.String() + "\nInput:\n" + m.Input + "\n"
	return u.String() == m.Answer, message
}

//-------------------
// CanonicalString -> Universe -> CanonicalString
// Tests: 1, 2
//-------------------
type canonicalToCanonical struct {
	Input, Answer string
}

func (c canonicalToCanonical) Check() (bool, string) {
	u := universe.LoadFromCanonicalString(c.Input)
	message := "Answer:\n" + c.Answer + "\nGuess:\n" + u.CanonicalString() + "\nInput:\n" + c.Input + "\n"
	return u.CanonicalString() == c.Answer, message
}

///
type testcase interface {
	Check() (bool, string)
}

type testsuite []testcase

func Run(t *testing.T, ts []testcase) {
	for _, test := range ts {
		ok, message := test.Check()
		if !ok {
			t.Errorf(message)
		}
	}
}

func TestCanonicalStringRoundtrip(t *testing.T) {
	u := universe.LoadFromCanonicalString("3,2,f|0,0,")
	if u.Width != 3 || u.Height != 2 {
		t.Error("Size does not match")
	}

	tests := []testcase{
		mapToNextMap{Input: "...\nOOO\n...\n", Answer: ".O.\n.O.\n.O.\n"},
		mapToNextMap{Input: ".O.\n.O.\n.O.\n", Answer: "...\nOOO\n...\n"},
		mapToNextMap{Input: "...\n...\n...\n", Answer: "...\n...\n...\n"},
		mapToNextMap{Input: "..\n..\n..\n", Answer: "..\n..\n..\n"},
		mapToNextMap{Input: ".OO.\n....\n.OO.\n", Answer: ".OO.\n....\n.OO.\n", Toroidal: true},
		mapToNextMap{Input: "O..O\n....\nO..O\n", Answer: "O..O\n....\nO..O\n", Toroidal: true},
		mapToNextMap{Input: "O..O\n....\nO..O\n", Answer: "....\n....\n....\n", Toroidal: false},
		canonicalToCanonical{Input: "3,2,true|0,0", Answer: "3,2,true|"},
		canonicalToCanonical{Input: "3,2,false|1,0", Answer: "3,2,false|"},
	}
	Run(t, tests)
}
