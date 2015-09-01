package goose

import "sync"

// ComputeArray represents an array based parallel computing architecture
type ComputeArray struct {
	elements [][]ComputeElement
}

// Ensure we implement ComputeElement for nesting purposes
var _ ComputeElement = (*ComputeArray)(nil)

// Return the layout of the elements of the simulation
func (arr ComputeArray) Shape() (rows, cols int) {
	els := arr.Elements()
	rows = len(els)
	cols = len(els[0])
	return
}

func (arr ComputeArray) Size() int {
	rows, cols := arr.Shape()
	return rows * cols
}

// Elements returns the ComputeElements which make up the array
func (arr *ComputeArray) Elements() [][]ComputeElement {
	return arr.elements
}

// Step causes the array to advance by single stepping each ComputeElement
func (arr *ComputeArray) Step() {
	var wg sync.WaitGroup
	wg.Add(arr.Size())
	for _, row := range arr.Elements() {
		for _, el := range row {
			go func(el ComputeElement) {
				defer wg.Done()
				el.Step()
			}(el)
		}
	}
	wg.Wait()
}

func (arr *ComputeArray) Data() [][]float64 {
	return nil
}

// Run causes the array to step as dictaded by the clk clock channel
func (arr *ComputeArray) Run(clk <-chan bool) {
	for range clk {
		arr.Step()
	}
}
