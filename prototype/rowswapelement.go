package goose

// RowSwapElement implements a ComputeElement which exchanges its top and bottom row
type RowSwapElement struct {
	data              [][]float64
	northIn, northOut chan []float64
	southIn, southOut chan []float64
}

func (el *RowSwapElement) Data() [][]float64 {
	return el.data
}

func (el *RowSwapElement) Swap() {
	sent := make(chan struct{})
	// Send - TODO benchmark with select
	go func() {
		el.northOut <- el.data[0]
		el.southOut <- el.data[len(el.data)-1]
		close(sent)
	}()

	var northRecv, southRecv []float64
	// Receive - TODO benchmark without select, just reversed ordering
	for i := 0; i < 2; i++ {
		select {
		case northRecv = <-el.northIn:
		case southRecv = <-el.southIn:
		}
	}

	// Update
	<-sent
	el.data[0] = northRecv
	el.data[len(el.data)-1] = southRecv
}

func (el *RowSwapElement) Step() {
	el.Swap()
}

func MakeRowSwapArray(els, elRows, elCols int) ComputeArray {
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
	return ComputeArray{elements: elems}
}
