package goose

import (
	"fmt"
	"sync"
)

// ComputeArray represents an array based parallel computing architecture
type ComputeArray struct {
	elements                   []ComputeElement
	elsVertical, elsHorizontal int // The number of row/cols of computeelements
}

// Ensure a ComputeArray is also a ComputeElement to support nesting
var _ ComputeElement = (*ComputeArray)(nil)

func NewComputeArray(elements []ComputeElement, elsVertical, elsHorizontal int) *ComputeArray {
	if len(elements) != elsVertical*elsHorizontal {
		panic(fmt.Sprintf("Unable to create %dx%d array from %d elements", elsHorizontal, elsVertical, len(elements)))
	}
	return &ComputeArray{elements, elsVertical, elsHorizontal}
}

// Data returns the amalgamated data held by the constituent compute elements
func (arr *ComputeArray) Data() [][]float64 {
	return [][]float64{[]float64{0}}
}

// Shape returns the (rows, cols) covered by our simulation domain
func (arr *ComputeArray) Shape() (int, int) {
	return len(arr.Data()), len(arr.Data()[0])
}

// Elements returns the ComputeElements which make up the array
func (arr ComputeArray) Elements() []ComputeElement {
	return arr.elements
}

// Step causes the array to advance by single stepping each ComputeElement
func (arr *ComputeArray) Step() {
	var wg sync.WaitGroup
	wg.Add(len(arr.Elements()))
	for _, el := range arr.Elements() {
		go func(el ComputeElement) {
			defer wg.Done()
			el.Step()
		}(el)
	}
	wg.Wait()
}

// Layout describes the number and layout of elements in this compute array
func (arr ComputeArray) Layout() (int, int) {
	return arr.elsVertical, arr.elsHorizontal
}

// Run causes the array to step as dictaded by the clk clock channel
func (arr *ComputeArray) Run(clk <-chan bool) {
	for range clk {
		arr.Step()
	}
}
