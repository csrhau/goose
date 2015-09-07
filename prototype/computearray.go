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
		panic(fmt.Sprintf("Unable to create %dx%d array from %d elements", elsVertical, elsHorizontal, len(elements)))
	}
	return &ComputeArray{elements, elsVertical, elsHorizontal}
}

// Data returns the amalgamated data held by the constituent compute elements
func (arr *ComputeArray) Data() [][]float64 {
	rows, _ := arr.Shape()
	ret := make([][]float64, rows)
	base := 0
	for vEl := 0; vEl < arr.elsVertical; vEl++ {
		elRows, _ := arr.ElementAt(vEl, 0).Shape()
		for hEl := 0; hEl < arr.elsHorizontal; hEl++ {
			el := arr.ElementAt(vEl, hEl)
			for rowOffset := 0; rowOffset < elRows; rowOffset++ {
				ret[base+rowOffset] = append(ret[base+rowOffset], el.Data()[rowOffset]...)
			}
		}
		base += elRows
	}
	return ret
}

// Shape returns the (rows, cols) covered by our simulation domain
func (arr *ComputeArray) Shape() (int, int) {
	rows, cols := 0, 0
	for i := 0; i < arr.elsVertical; i++ {
		r, _ := arr.ElementAt(i, 0).Shape()
		rows += r
	}
	for i := 0; i < arr.elsHorizontal; i++ {
		_, c := arr.ElementAt(0, i).Shape()
		cols += c
	}
	return rows, cols
}

// Layout describes the number and layout of elements in this compute array
func (arr ComputeArray) Layout() (int, int) {
	return arr.elsVertical, arr.elsHorizontal
}

// ElementAt allows indexing into the sub-elements of the array
func (arr ComputeArray) ElementAt(i, j int) ComputeElement {
	return arr.Elements()[i*arr.elsHorizontal+j]
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

// Run causes the array to step as dictaded by the clk clock channel
func (arr *ComputeArray) Run(clk <-chan bool) {
	for range clk {
		arr.Step()
	}
}
