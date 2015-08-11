package goose

import "testing"

func TestCountingElementCycles(t *testing.T) {
	testCycles := 1337
	c := make(chan bool)
	d := make([][]float64, 0)
	el := CountingElement{0, c, d}
	go el.Step()
	for i := 0; i < testCycles; i++ {
		c <- true
	}
	close(c)
	cycles := el.Cycles()
	if cycles != testCycles {
		t.Error("Cycle Count Mismatch, expected", testCycles, "got", cycles)
	}
}
