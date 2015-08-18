package goose

// CountingElement is a dummy ComputeElement which simply counts cycles
type CountingElement struct {
	data       [][]float64
	iterations int
}

// Data returns the data stored by this element
func (el *CountingElement) Data() [][]float64 {
	return el.data
}

// Step causes this element to advance by one step, setting each cell of this
// element's data to the step count
func (el *CountingElement) Step() {
	el.iterations++
	for _, r := range el.data {
		for i := range r {
			r[i] = float64(el.iterations)
		}
	}
}

// Iterations returns the number of times this element has stepped for testing
func (el CountingElement) Iterations() int {
	return el.iterations
}
