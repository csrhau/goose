package goose

import "testing"

func TestRowSwapElementSwapsRows(t *testing.T) {
	northIn := make(chan []float64)
	northOut := make(chan []float64)
	southIn := make(chan []float64)
	southOut := make(chan []float64)
	initData := [][]float64{
		[]float64{1.1, 1.2, 1.3},
		[]float64{2.1, 2.2, 2.3},
		[]float64{3.1, 3.2, 3.3},
	}
	ni := []float64{0.1, 0.2, 0.3}
	si := []float64{4.1, 4.2, 4.3}

	el := new(RowSwapElement)
	el.data = initData
	el.northIn = northIn
	el.northOut = northOut
	el.southIn = southIn
	el.southOut = southOut

	// Make sure there is something waiting to be swapped
	go func() {
		northIn <- ni
		southIn <- si
	}()

	// Run an iteration of the element logic
	go el.Step()

	for i := 0; i < 2; i++ {
		select {
		case no := <-northOut:
			for i, v := range no {
				if v != initData[0][i] {
					t.Error("mismatch in received data! expected", v, "got", initData[0][i])
				}
			}
		case so := <-southOut:
			for i, v := range so {
				if v != initData[2][i] {
					t.Error("mismatch in received data! expected", v, "got", initData[2][i])
				}
			}
		}
	}
}

func TestRowSwapElementInArray(t *testing.T) {
	const els = 42
	const elRows = 5
	const elCols = 1000
	elems := make([]ComputeElement, els)
	upChans := make([]chan []float64, els)
	downChans := make([]chan []float64, els)

	for i := 0; i < els; i++ {
		upChans[i] = make(chan []float64)
		downChans[i] = make(chan []float64)
	}

	for i := 0; i < els; i++ {
		data := make([][]float64, elRows)
		for j := 0; j < elRows; j++ {
			data[j] = make([]float64, elCols)
			// Populate each element with the float representation of it's index
			for k := 0; k < elCols; k++ {
				data[j][k] = float64(i)
			}
		}
		elem := RowSwapElement{
			data:     data,
			northIn:  downChans[i],
			northOut: upChans[i],
			southIn:  upChans[(i+1)%els], // Wrap around boundary
			southOut: downChans[(i+1)%els],
		}
		elems[i] = &elem
	}
	ar := new(ComputeArray)
	ar.elements = elems

	const testIters = 10

	for i := 0; i < 10; i++ {
		ar.Step()

	}
}