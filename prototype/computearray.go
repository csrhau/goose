package goose

import "sync"

// ComputeArray represents an array based parallel computing architecture
type ComputeArray struct {
	elements []ComputeElement
}

// Ensure we implement ComputeElement for nesting purposes
var _ ComputeElement = (*ComputeArray)(nil)

// Elements returns the ComputeElements which make up the array
func (arr *ComputeArray) Elements() []ComputeElement {
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

func (arr *ComputeArray) Data() [][]float64 {

	return nil
}

// Run causes the array to step as dictaded by the clk clock channel
func (arr *ComputeArray) Run(clk <-chan bool) {
	for range clk {
		arr.Step()
	}
}
