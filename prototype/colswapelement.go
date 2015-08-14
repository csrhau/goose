package goose

// ColSwapElement implements a ComputeElement which exchanges its left and right cols
type ColSwapElement struct {
	rows, cols      int
	data            [][]float64
	westIn, westOut chan []float64
	eastIn, eastOut chan []float64
}

func (el *ColSwapElement) Data() [][]float64 {
	return el.data
}

func (el *ColSwapElement) Swap() {
	// Send
	go func() {
		westOutBuff := make([]float64, el.rows)
		eastOutBuff := make([]float64, el.rows)
		for r := 0; r < el.rows; r++ {
			westOutBuff[r] = el.data[r][1]
			eastOutBuff[r] = el.data[r][el.cols-2]
		}
		el.westOut <- westOutBuff
		el.eastOut <- eastOutBuff
	}()

	// Receive
	eastInBuff := <-el.eastIn
	westInBuff := <-el.westIn
	for r := 0; r < el.rows; r++ {
		el.data[r][el.cols-1] = eastInBuff[r]
		el.data[r][0] = westInBuff[r]
	}
}

func (el *ColSwapElement) Step() {
	el.Swap()
}

func MakeColSwapArray(els, elRows, elCols int) ComputeArray {
	elems := make([]ComputeElement, els)
	firstRight, firstLeft := make(chan []float64), make(chan []float64)
	lastLeft, lastRight := firstLeft, firstRight
	for i := 0; i < els; i++ {
		data := make([][]float64, elRows)
		for j := 0; j < elRows; j++ {
			data[j] = make([]float64, elCols)
			for k := 0; k < elCols; k++ {
				data[j][k] = float64(i)
			}
		}
		rse := new(ColSwapElement)
		rse.rows = elRows
		rse.cols = elCols
		rse.data = data
		rse.westIn, rse.westOut = lastRight, lastLeft
		if i < els-1 {
			rse.eastIn, rse.eastOut = make(chan []float64), make(chan []float64)
		} else {
			rse.eastIn, rse.eastOut = firstLeft, firstRight
		}
		lastRight, lastLeft = rse.eastOut, rse.eastIn
		elems[i] = rse
	}
	return ComputeArray{elements: elems}
}
