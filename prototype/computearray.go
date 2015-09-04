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

func NewComputeArray(elements []ComputeElement, elsVertical, elsHorizontal int) *ComputeArray {
	if len(elements) != elsVertical*elsHorizontal {
		panic(fmt.Sprintf("Unable to create %dx%d array from %d elements", elsHorizontal, elsVertical, len(elements)))
	}
	return &ComputeArray{elements, elsVertical, elsHorizontal}
}

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

// Layout
func (arr ComputeArray) Layout() (int, int) {
	return arr.elsVertical, arr.elsHorizontal
}

// Run causes the array to step as dictaded by the clk clock channel
func (arr *ComputeArray) Run(clk <-chan bool) {
	for range clk {
		arr.Step()
	}
}
