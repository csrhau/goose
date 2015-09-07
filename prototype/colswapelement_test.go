package goose

import "testing"

func ColSwapElementIsAComputeElement(t *testing.T) {
	var el ComputeElement
	el = new(ColSwapElement)
	_ = el
}

func TestColSwapElementSwapsCols(t *testing.T) {
	westIn := make(chan []float64)
	westOut := make(chan []float64)
	eastIn := make(chan []float64)
	eastOut := make(chan []float64)
	initData := [][]float64{
		[]float64{1.1, 1.2, 1.3, 1.4},
		[]float64{2.1, 2.2, 2.3, 2.4},
		[]float64{3.1, 3.2, 3.3, 3.4},
		[]float64{4.1, 4.2, 4.3, 4.4},
	}
	wi := []float64{1.0, 2.0, 3.0, 4.0}
	ei := []float64{1.5, 2.5, 3.5, 4.5}

	el := new(ColSwapElement)
	el.rows = 4
	el.cols = 4
	el.data = initData
	el.westIn = westIn
	el.westOut = westOut
	el.eastIn = eastIn
	el.eastOut = eastOut

	// Make sure there is something waiting to be swapped
	go func() {
		westIn <- wi
		eastIn <- ei
	}()

	// Run an iteration of the element logic
	go el.Step()

	// See what comes out
	for i := 0; i < 2; i++ {
		select {
		case wo := <-westOut:
			if len(wo) != 4 {
				t.Error("Received", len(wo), "elements from wo, expected 4")
			}
			for i, v := range wo {
				if v != initData[i][1] {
					t.Error("mismatch in west received data! got", v, "expected", initData[i][1])

				}
			}
		case eo := <-eastOut:
			if len(eo) != 4 {
				t.Error("Received", len(eo), "elements from eo, expected 4")
			}
			for i, v := range eo {
				if v != initData[i][2] {
					t.Error("mismatch in received data! got", v, "expected", initData[i][2])
				}
			}
		}
	}
}

func TestColSwapElementInArrayMakeConstruction(t *testing.T) {
	const els = 3
	const elRows = 4
	const elCols = 4
	arr := MakeColSwapArray(els, elRows, elCols)
	// Expect Contiguous numbers
	for i, e := range arr.elements {
		for j := 0; j < elRows; j++ {
			for k := 0; k < elCols; k++ {
				if e.Data()[j][k] != float64(i) {
					t.Error("Mismatched data on unshuffled step")
				}
			}
		}
	}
	// Expect Contiguous numbers
	for i, e := range arr.elements {
		for j := 0; j < elRows; j++ {
			for k := 0; k < elCols; k++ {
				if e.Data()[j][k] != float64(i) {
					t.Error("Mismatched data on unshuffled step")
				}
			}
		}
	}
	for s := 0; s < 10; s++ {
		arr.Step()
		for i, e := range arr.elements {
			for k := 0; k < elRows; k++ {
				prev, next := i-1, i+1
				if prev < 0 {
					prev = els - 1
				}
				if next > els-1 {
					next = 0
				}
				// West Row from prev. elements
				if e.Data()[k][0] != float64(prev) {
					t.Error("Mismatch data on top row of shuffled step")
				}
				// Inner rows from self
				for j := 1; j < 3; j++ {
					if e.Data()[k][j] != float64(i) {
						t.Error("Mismatch in inner unshuffled data")
					}
				}

				// East Row from next element
				if e.Data()[k][3] != float64(next) {
					t.Error("Mismatch data on bottom row of shuffled step")
				}
			}
		}
	}
}

func TestColSwapElementShapeWorksAsExpected(t *testing.T) {
	var arr *ComputeArray
	for elsHorizontal := 1; elsHorizontal < 4; elsHorizontal++ {
		for elRows := 1; elRows < 10; elRows++ {
			for elCols := 1; elCols < 10; elCols++ {
				arr = MakeColSwapArray(elsHorizontal, elRows, elCols)
				for _, el := range arr.Elements() {
					er, ec := el.Shape()
					if er != elRows {
						t.Error("Misshapen element rows, expected", elRows, "got", er)
					}
					if ec != elCols {
						t.Error("Misshapen element cols, expected", elCols, "got", ec)
					}
					dat := el.Data()

					if len(dat) != er {
						t.Error("Element reporting inconsistent number of rows!")
					}
					if len(dat[0]) != ec {
						t.Error("Element reporting inconsistent number of cols!")
					}
				}
			}
		}
	}
}

func BenchmarkColSwap(b *testing.B) {
	arr := MakeColSwapArray(25*25, 100, 100)
	for i := 0; i < b.N; i++ {
		arr.Step()
	}
}
