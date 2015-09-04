package goose

import "testing"

func TestCountingElementIterations(t *testing.T) {
	el := new(CountingElement) //NB this is done so el is a pointer!
	if el.Iterations() != 0 {
		t.Error("element initialized with non-zero stepcount")
	}
	for i := 1; i < 10; i++ {
		el.Step()
		if el.Iterations() != i {
			t.Error("Iteration count mismatch: ", el.Iterations(), i)
		}
	}
}

func TestCountingElementData(t *testing.T) {
	el := new(CountingElement)
	el.data = [][]float64{
		make([]float64, 3),
		make([]float64, 3),
		make([]float64, 3),
	}
	for i := 1; i < 10; i++ {
		el.Step()
		for _, r := range el.data {
			for _, v := range r {
				if v != float64(i) {
					t.Error("Counting element data error")
				}
			}
		}
	}
}

func TestCountingElementInArray(t *testing.T) {
	el := new(CountingElement)
	el.data = [][]float64{
		make([]float64, 3),
		make([]float64, 3),
		make([]float64, 3),
	}
	ar := NewComputeArray([]ComputeElement{el}, 1, 1)
	for i := 1; i < 10; i++ {
		ar.Step()
		el := ar.Elements()[0]
		for _, r := range el.Data() {
			for _, v := range r {
				if v != float64(i) {
					t.Error("Counting element data error")
				}
			}
		}
	}
}
