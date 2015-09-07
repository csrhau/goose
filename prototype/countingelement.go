package goose

// CountingElement is a dummy ComputeElement which simply counts cycles
type CountingElement struct {
	data    [][]float64
	counter int
}

// Data returns the data stored by this element
func (el *CountingElement) Data() [][]float64 {
	return el.data
}

// Shape returns the (rows, cols) covered by our simulation domain
func (el *CountingElement) Shape() (int, int) {
	return len(el.Data()), len(el.Data()[0])
}

// Step causes this element to advance by one step, setting each cell of this
// element's data to the step count
func (el *CountingElement) Step() {
	el.counter++
	for _, r := range el.data {
		for i := range r {
			r[i] = float64(el.counter)
		}
	}
}

// Iterations returns the number of times this element has stepped for testing
func (el CountingElement) Iterations() int {
	return el.counter
}

func MakeCountingArray(elsVertical, elsHorizontal, elRows, elCols int) *ComputeArray {
	els := elsVertical * elsHorizontal
	elems := make([]ComputeElement, els)
	for e := 0; e < els; e++ {
		data := make([][]float64, elRows)
		for i := 0; i < elRows; i++ {
			data[i] = make([]float64, elCols)
			for j := 0; j < elCols; j++ {
				data[i][j] = float64(e)
			}
		}
		elems[e] = &CountingElement{data, e}
	}
	return NewComputeArray(elems, elsVertical, elsHorizontal)
}
