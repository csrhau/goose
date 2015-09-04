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

func TestMakeCountingArrayPopulatesData(t *testing.T) {
	elRows := 5
	elCols := 7
	for elsHorizontal := 1; elsHorizontal < 8; elsHorizontal++ {
		for elsVertical := 1; elsVertical < 8; elsVertical++ {
			arr := MakeCountingArray(elsVertical, elsHorizontal, elRows, elCols)
			for i, el := range arr.Elements() {
				if len(el.Data()) != elRows || len(el.Data()[0]) != elCols {
					t.Error("Misshapen internal data detected:", len(el.Data()), "x", len(el.Data()[0]))
				}
				for _, r := range el.Data() {
					for _, c := range r {
						if c != float64(i) {
							t.Error("Incorrect data in initialization, expected", float64(i), "got", c, "for el", i)
						}
					}
				}
			}
		}
	}
}

func TestMakeCountingArrayProducesViableArray(t *testing.T) {
	elRows := 71
	elCols := 2
	elsHorizontal := 13
	elsVertical := 15
	arr := MakeCountingArray(elsVertical, elsHorizontal, elRows, elCols)
	for step := 0; step < 10; step++ {
		for i, el := range arr.Elements() {
			for _, r := range el.Data() {
				for _, c := range r {
					if c != float64(i+step) {
						t.Error("Incorrect data in step", step, "expected", float64(i), "got", c, "for el", i)
					}
				}
			}
		}
		arr.Step()
	}
}
