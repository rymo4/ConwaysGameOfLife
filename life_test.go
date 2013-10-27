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
	u.Next()
	next := u.CanonicalString()
	answer := "3,2|"
	if next != answer {
		t.Error("Canonical strings do not match: %s != %s", next, answer)
	}
}
