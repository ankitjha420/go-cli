package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("one two three\tfour five\n")
	ans := 5

	res := count(b, false)
	if res != ans {
		t.Errorf("Expected %d but got %d", ans, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("one\n two\n three")
	ans := 3

	res := count(b, true)
	if res != ans {
		t.Errorf("Expected %d but got %d", ans, res)
	}
}
