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
		for _, el := range ar.Elements() {
			for _, r := range el.Data() {
				for _, v := range r {
					if v != float64(i) {
						t.Error("Counting element data error")
					}
				}
			}
		}
	}
}

func TestArrayClocking(t *testing.T) {
	ar := new(ComputeArray)
	ar.elements = []ComputeElement{
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
	}
	clk := make(chan bool)
	go ar.Run(clk)
	for i := 1; i < 10; i++ {
		clk <- true
		for _, el := range ar.Elements() {
			for _, r := range el.Data() {
				for _, v := range r {
					if v != float64(i) {
						t.Error("Counting element data error")
					}
				}
			}
		}
	}
	close(clk)
}

func TestArrayData(t *testing.T) {
	ar := new(ComputeArray)
	ar.elements = []ComputeElement{
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
	}
	clk := make(chan bool)
	go ar.Run(clk)
	for i := 1; i < 10; i++ {
		if ar.Data() == nil {
			t.Error("Nil data returned!")

		}
	}
	close(clk)
}
