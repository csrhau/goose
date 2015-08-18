package goose

// RowSwapElement implements a ComputeElement which exchanges its top and bottom row
type RowSwapElement struct {
	data              [][]float64
	northIn, northOut chan []float64
	southIn, southOut chan []float64
}

// Data returns the data stored by this element
func (el *RowSwapElement) Data() [][]float64 {
	return el.data
}

// Swap causes this element to exchange data with its neighbours
func (el *RowSwapElement) Swap() {
	// Send
	go func() {
		el.northOut <- el.data[1]
		el.southOut <- el.data[len(el.data)-2]
	}()
	// Receive
	el.data[len(el.data)-1] = <-el.southIn
	el.data[0] = <-el.northIn
}

// Step causes this element to advance by one step
func (el *RowSwapElement) Step() {
	el.Swap()
}

// MakeRowSwapArray constructs a ComputeArray populated by RowSwapElements
func MakeRowSwapArray(els, elRows, elCols int) *ComputeArray {
	elems := make([]ComputeElement, els)
	topDown, topUp := make(chan []float64), make(chan []float64)
	lastUp, lastDown := topUp, topDown
	for i := 0; i < els; i++ {
		data := make([][]float64, elRows)
		for j := 0; j < elRows; j++ {
			data[j] = make([]float64, elCols)
			for k := 0; k < elCols; k++ {
				data[j][k] = float64(i)
			}
		}
		rse := new(RowSwapElement)
		rse.data = data
		rse.northIn, rse.northOut = lastDown, lastUp
		if i < els-1 {
			rse.southIn, rse.southOut = make(chan []float64), make(chan []float64)
		} else {
			rse.southIn, rse.southOut = topUp, topDown
		}
		lastDown, lastUp = rse.southOut, rse.southIn
		elems[i] = rse
	}
	return &ComputeArray{elements: elems}
}
