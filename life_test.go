package main

import (
	"github.com/rymo4/life/universe"
	"testing"
)

func TestCanonicalStringRoundtrip(t *testing.T) {
	u := universe.LoadFromCanonicalString("3,2|0,0,")
	if u.Width != 3 || u.Height != 2 {
		t.Error("Size does not match")
	}

  runRoundtrip("3,2|0,0,", "3,2|", t)
  runRoundtrip("3,3|1,0,1,1,1,2,", "3,3|0,1,1,1,2,1,", t)

}

func runRoundtrip(first, answer string, t *testing.T) {
	u := universe.LoadFromCanonicalString(first)
	u.Next()
	next := u.CanonicalString()
	if next != answer {
		t.Error("Canonical strings do not match: %s != %s", next, answer)
	}
}
