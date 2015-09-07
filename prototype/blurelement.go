package goose

import "fmt"

// BlurElement implements a ComputeElement which perfoms a gaussian blur
type BlurElement struct {
	rows, cols        int
	data, scratch     [][]float64
	northIn, northOut chan []float64
	westIn, westOut   chan []float64
	southIn, southOut chan []float64
	eastIn, eastOut   chan []float64
}

// Ensure we implement ComputeElement
var _ ComputeElement = (*BlurElement)(nil)

// Data returns the data stored by this element
func (el *BlurElement) Data() [][]float64 {
	rows, cols := el.Shape()
	data := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		data[i] = make([]float64, cols)
	}
	return data
}

// Shape returns the (rows, cols) covered by our simulation domain
func (el BlurElement) Shape() (int, int) {
	return el.rows, el.cols
}

// Swap causes this element to exchange data with its neighbours
func (el *BlurElement) Swap() {
	// Send
	go func() {
		westOutBuff := make([]float64, el.rows)
		eastOutBuff := make([]float64, el.rows)
		for r := 0; r < el.rows; r++ {
			westOutBuff[r] = el.data[r+1][1]
			fmt.Println("coming acropper accessing", r+1, el.cols, "of", len(el.data), len(el.data[0]))
			fmt.Println("el.rows:", el.rows, "cols:", el.cols)
			eastOutBuff[r] = el.data[r+1][el.cols]
		}
		el.northOut <- el.data[1][1 : el.cols+1]
		el.southOut <- el.data[el.rows][1 : el.cols+1]
		el.westOut <- westOutBuff
		el.eastOut <- eastOutBuff
	}()

	// Receive
	southInBuff := <-el.southIn
	northInBuff := <-el.northIn
	eastInBuff := <-el.eastIn
	westInBuff := <-el.westIn
	for c := 0; c < el.cols; c++ {
		el.data[0][c+1] = northInBuff[c]
		el.data[el.rows+1][c+1] = southInBuff[c]
	}
	for r := 0; r < el.rows; r++ {
		el.data[r+1][el.cols+1] = eastInBuff[r]
		el.data[r+1][0] = westInBuff[r]
	}
}

// Step causes this element to advance by one step
func (el *BlurElement) Step() {
	el.Swap()
	// Gaussian Blur
	kernel := [][]float64{
		[]float64{1 / 16.0, 1 / 8.0, 1 / 16.0},
		[]float64{1 / 8.0, 1 / 4.0, 1 / 8.0},
		[]float64{1 / 16.0, 1 / 8.0, 1 / 16.0},
	}
	for i := 1; i < el.rows+1; i++ {
		for j := 1; j < el.cols+1; j++ {
			el.scratch[i][j] = el.data[i-1][j-1]*kernel[0][0] +
				el.data[i-1][j]*kernel[0][1] +
				el.data[i-1][j+1]*kernel[0][2] +
				el.data[i][j-1]*kernel[1][0] +
				el.data[i][j]*kernel[1][1] +
				el.data[i][j+1]*kernel[1][2] +
				el.data[i+1][j-1]*kernel[2][0] +
				el.data[i+1][j]*kernel[2][1] +
				el.data[i+1][j+1]*kernel[2][2]

		}
	}
	el.data, el.scratch = el.scratch, el.data
}

// MakeBlurArray constructs a ComputeArray populated by BlurElements
func MakeBlurArray(data [][]float64, widthEls, heightEls int) ComputeArray {
	els := widthEls * heightEls

	if len(data)%heightEls != 0 {
		panic("Imbalanced cell distribution")
	}
	elRows := len(data) / heightEls
	elCols := len(data[0]) / widthEls

	cseelems := make([]*BlurElement, els)

	for ely := 0; ely < heightEls; ely++ {
		for elx := 0; elx < widthEls; elx++ {
			elID := ely*widthEls + elx
			elData := make([][]float64, elRows+2)
			scratch := make([][]float64, elRows+2)
			for i := 0; i < elRows+2; i++ {
				elData[i] = make([]float64, elCols+2)
				scratch[i] = make([]float64, elCols+2)
			}

			for i := 1; i < elRows+1; i++ {
				for j := 1; j < elCols+2-1; j++ {
					yOff := ely * elRows
					xOff := elx * elCols
					rowAccess := i + yOff - 1
					colAccess := j + xOff - 1
					elData[i][j] = data[rowAccess][colAccess]
				}
			}

			cse := new(BlurElement)
			cse.rows = elRows
			cse.cols = elCols
			cse.data = elData
			cse.scratch = scratch
			cse.northOut = make(chan []float64)
			cse.westOut = make(chan []float64)
			cse.southOut = make(chan []float64)
			cse.eastOut = make(chan []float64)
			cseelems[elID] = cse
		}
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

	elems := make([]ComputeElement, len(cseelems))
	for i, e := range cseelems {
		elems[i] = e
	}

	return ComputeArray{elements: elems}
}
