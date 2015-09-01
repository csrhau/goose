package goose

// CartesianSwapElement implements a ComputeElement which exchanges its left and right cols
type CartesianSwapElement struct {
	rows, cols        int
	data              [][]float64
	northIn, northOut chan []float64
	westIn, westOut   chan []float64
	southIn, southOut chan []float64
	eastIn, eastOut   chan []float64
}

// Ensure we implement ComputeElement
var _ ComputeElement = (*CartesianSwapElement)(nil)

// Data returns the data stored by this element
func (el *CartesianSwapElement) Data() [][]float64 {
	return el.data
}

// Swap causes this element to exchange data with its neighbours
func (el *CartesianSwapElement) Swap() {
	// Send
	go func() {
		westOutBuff := make([]float64, el.rows-2)
		eastOutBuff := make([]float64, el.rows-2)
		for r := 0; r < el.rows-2; r++ {
			westOutBuff[r] = el.data[r+1][1]
			eastOutBuff[r] = el.data[r+1][el.cols-2]
		}
		el.northOut <- el.data[1][1 : el.cols-1]
		el.southOut <- el.data[el.rows-2][1 : el.cols-1]
		el.westOut <- westOutBuff
		el.eastOut <- eastOutBuff
	}()

	// Receive
	southInBuff := <-el.southIn
	northInBuff := <-el.northIn
	eastInBuff := <-el.eastIn
	westInBuff := <-el.westIn
	for c := 0; c < el.cols-2; c++ {
		el.data[0][c+1] = northInBuff[c]
		el.data[el.rows-1][c+1] = southInBuff[c]
	}
	for r := 0; r < el.rows-2; r++ {
		el.data[r+1][el.cols-1] = eastInBuff[r]
		el.data[r+1][0] = westInBuff[r]
	}
}

// Step causes this element to advance by one step
func (el *CartesianSwapElement) Step() {
	el.Swap()
}

// MakeCartesianSwapArray constructs a ComputeArray populated by CartesianSwapElements
func MakeCartesianSwapArray(widthEls, heightEls, elRows, elCols int) ComputeArray {
	els := widthEls * heightEls
	cseelems := make([]*CartesianSwapElement, els)
	for e := 0; e < els; e++ {
		data := make([][]float64, elRows)
		for i := 0; i < elRows; i++ {
			data[i] = make([]float64, elCols)
			for j := 0; j < elCols; j++ {
				data[i][j] = float64(e)
			}
		}
		cse := new(CartesianSwapElement)
		cse.rows = elRows
		cse.cols = elCols
		cse.data = data
		cse.northOut = make(chan []float64)
		cse.westOut = make(chan []float64)
		cse.southOut = make(chan []float64)
		cse.eastOut = make(chan []float64)
		cseelems[e] = cse
	}

	// Weird naming convention, but it is correct
	for row := 0; row < heightEls; row++ {
		for col := 0; col < widthEls; col++ {
			// Find neighbours, dealing with wrap-around
			nextRow := (row + 1) % heightEls
			prevRow := (row - 1 + heightEls) % heightEls
			nextCol := (col + 1) % widthEls
			prevCol := (col - 1 + widthEls) % widthEls
			cse := cseelems[widthEls*row+col]
			cse.northIn = cseelems[widthEls*prevRow+col].southOut
			cse.westIn = cseelems[widthEls*row+prevCol].eastOut
			cse.southIn = cseelems[widthEls*nextRow+col].northOut
			cse.eastIn = cseelems[widthEls*row+nextCol].westOut
		}
	}

	elems := make([]ComputeElement, els)
	for i, v := range cseelems {
		elems[i] = v
	}
	return ComputeArray{elements: elems}
}
