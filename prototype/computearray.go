package goose

import "sync"

type ComputeArray struct {
	elements []ComputeElement
}

func (arr *ComputeArray) Elements() []ComputeElement {
	return arr.elements
}

func (arr *ComputeArray) Step() {
	var wg sync.WaitGroup
	wg.Add(len(arr.Elements()))
	for _, el := range arr.Elements() {
		go func() {
			defer wg.Done()
			el.Step()
		}()
	}
	wg.Wait()
}

func (arr *ComputeArray) Run(clk <-chan bool) {
	for range clk {
		arr.Step()
	}
}
