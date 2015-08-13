package goose

// CountingElement is a dummy ComputeElement which simply counts cycles
type CountingElement struct {
	data       [][]float64
	iterations int
}

func (el *CountingElement) Data() [][]float64 {
	return el.data
}

func (el *CountingElement) Step() {
	el.iterations++
	for _, r := range el.data {
		for i := range r {
			r[i] = float64(el.iterations)
		}
	}
}

func (el CountingElement) Iterations() int {
	return el.iterations
}
