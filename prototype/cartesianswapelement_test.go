package goose

import "testing"

func CartesianSwapElementIsAComputeElement(t *testing.T) {
	var el ComputeElement
	el = new(CartesianSwapElement)
	_ = el
}

func TestCartesianSwapElementSwapsCols(t *testing.T) {
	northIn := make(chan []float64)
	northOut := make(chan []float64)
	westIn := make(chan []float64)
	westOut := make(chan []float64)
	southIn := make(chan []float64)
	southOut := make(chan []float64)
	eastIn := make(chan []float64)
	eastOut := make(chan []float64)
	initData := [][]float64{
		[]float64{1.1, 1.2, 1.3, 1.4},
		[]float64{2.1, 2.2, 2.3, 2.4},
		[]float64{3.1, 3.2, 3.3, 3.4},
		[]float64{4.1, 4.2, 4.3, 4.4},
	}

	ni := []float64{0.0, 0.0}
	wi := []float64{0.1, 0.1}
	si := []float64{0.2, 0.2}
	ei := []float64{0.3, 0.3}

	el := new(CartesianSwapElement)
	el.rows = 4
	el.cols = 4
	el.data = initData
	el.northIn = northIn
	el.northOut = northOut
	el.westIn = westIn
	el.westOut = westOut
	el.southIn = southIn
	el.southOut = southOut
	el.eastIn = eastIn
	el.eastOut = eastOut

	// Make sure there is something waiting to be swapped
	go func() {
		southIn <- si
		northIn <- ni
		eastIn <- ei
		westIn <- wi
	}()

	// Run an iteration of the element logic
	go el.Step()

	// See what comes out
	for i := 0; i < 4; i++ {
		select {
		case no := <-northOut:
			if len(no) != len(ni) {
				t.Error("Recieved", len(no), "elements from no, expected", len(ni))
			}
			for i, v := range no {
				exp := initData[1][i+1]
				if v != exp {
					t.Error("mismatch in west received data! got", v, "expected", exp)
				}
			}
		case wo := <-westOut:
			if len(wo) != len(wi) {
				t.Error("Recieved", len(wo), "elements from wo, expected", len(wi))
			}
			for i, v := range wo {
				exp := initData[i+1][1]
				if v != exp {
					t.Error("mismatch in west received data! got", v, "expected", exp)
				}
			}
		case so := <-southOut:
			if len(so) != len(si) {
				t.Error("Recieved", len(so), "elements from so, expected", len(si))
			}
			for i, v := range so {
				exp := initData[2][i+1]
				if v != exp {
					t.Error("mismatch in west received data! got", v, "expected", exp)
				}
			}
		case eo := <-eastOut:
			if len(eo) != len(ei) {
				t.Error("Recieved", len(eo), "elements from eo, expected", len(ei))
			}
			for i, v := range eo {
				exp := initData[i+1][2]
				if v != exp {
					t.Error("mismatch in received data! got", v, "expected", exp)
				}
			}
		}
	}
}

func TestMakeCartesianSwapArrayPopulatesData(t *testing.T) {
	elRows := 5
	elCols := 7
	for widthEls := 1; widthEls < 8; widthEls++ {
		for heightEls := 1; heightEls < 8; heightEls++ {
			arr := MakeCartesianSwapArray(widthEls, heightEls, elRows, elCols)
			for i, el := range arr.Elements() {
				if len(el.Data()) != elRows || len(el.Data()[0]) != elCols {
					t.Error("Misshapen internal data detected:", len(el.Data()), "x", len(el.Data()[0]))
				}
				for _, r := range el.Data() {
					for _, c := range r {
						if c != float64(i) {
							t.Error("Incorrect data in initialization, expected", float64(i), "got", c, "for el", i)
						}
					}
				}
			}
		}
	}
}

func TestCartesianSwapArraySwapsSuccessfully(t *testing.T) {
	arr := MakeCartesianSwapArray(2, 2, 4, 4)

	initialExpected := [][]float64{
		[]float64{0, 0, 0, 0, 1, 1, 1, 1},
		[]float64{0, 0, 0, 0, 1, 1, 1, 1},
		[]float64{0, 0, 0, 0, 1, 1, 1, 1},
		[]float64{0, 0, 0, 0, 1, 1, 1, 1},
		[]float64{2, 2, 2, 2, 3, 3, 3, 3},
		[]float64{2, 2, 2, 2, 3, 3, 3, 3},
		[]float64{2, 2, 2, 2, 3, 3, 3, 3},
		[]float64{2, 2, 2, 2, 3, 3, 3, 3},
	}

	initialReceived := make([][]float64, 8)
	elsTmp := arr.Elements()

	for i := 0; i < 4; i++ {
		initialReceived[i] = append(elsTmp[0].Data()[i], elsTmp[1].Data()[i]...)
		initialReceived[i+4] = append(elsTmp[2].Data()[i], elsTmp[3].Data()[i]...)
	}

	for i, r := range initialReceived {
		for j, c := range r {
			if initialExpected[i][j] != c {
				t.Error("Mismatched initialization data, expected", initialExpected[i][j], "got", c)
			}
		}
	}

	steppedExpected := [][]float64{
		[]float64{0, 2, 2, 0, 1, 3, 3, 1},
		[]float64{1, 0, 0, 1, 0, 1, 1, 0},
		[]float64{1, 0, 0, 1, 0, 1, 1, 0},
		[]float64{0, 2, 2, 0, 1, 3, 3, 1},
		[]float64{2, 0, 0, 2, 3, 1, 1, 3},
		[]float64{3, 2, 2, 3, 2, 3, 3, 2},
		[]float64{3, 2, 2, 3, 2, 3, 3, 2},
		[]float64{2, 0, 0, 2, 3, 1, 1, 3},
	}

	for step := 0; step < 5; step++ {
		arr.Step()
		steppedReceived := make([][]float64, 8)
		for i := 0; i < 4; i++ {
			steppedReceived[i] = append(elsTmp[0].Data()[i], elsTmp[1].Data()[i]...)
			steppedReceived[i+4] = append(elsTmp[2].Data()[i], elsTmp[3].Data()[i]...)
		}

		for i, r := range steppedReceived {
			for j, c := range r {
				if steppedExpected[i][j] != c {
					t.Error("Mismatched steppedization data, expected", steppedExpected[i][j], "got", c)
				}
			}
		}
	}
}

func TestMakeCartesianSwapArrayReturnsViableArray(t *testing.T) {
	arr := MakeCartesianSwapArray(3, 5, 10, 10)
	arr.Step()
}

func BenchmarkCartesianSwap(b *testing.B) {
	arr := MakeCartesianSwapArray(25, 25, 100, 100)
	for i := 0; i < b.N; i++ {
		arr.Step()
	}
}
