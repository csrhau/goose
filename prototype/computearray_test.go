package goose

import "testing"

func TestArrayIterations(t *testing.T) {
	ar := new(ComputeArray)
	ar.elements = []ComputeElement{
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
	}
	for i := 1; i < 10; i++ {
		ar.Step()
	}
}
