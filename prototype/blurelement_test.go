package goose

import "testing"

func BlurElementIsAComputeElement(t *testing.T) {
	var el ComputeElement
	el = new(BlurElement)
	_ = el
}

func TestBlurElementBlurs(t *testing.T) {
	northIn := make(chan []float64)
	northOut := make(chan []float64)
	westIn := make(chan []float64)
	westOut := make(chan []float64)
	southIn := make(chan []float64)
	southOut := make(chan []float64)
	eastIn := make(chan []float64)
	eastOut := make(chan []float64)
	initData := [][]float64{
		[]float64{0, 0, 0, 0, 0},
		[]float64{0, 1, 1, 1, 0},
		[]float64{0, 1, 2, 1, 0},
		[]float64{0, 1, 1, 1, 0},
		[]float64{0, 0, 0, 0, 0},
	}
	scratch := make([][]float64, 5)
	for i := 0; i < 5; i++ {
		scratch[i] = make([]float64, 5)
	}

	ni := []float64{0.0, 0.0, 0.0}
	wi := []float64{0.0, 0.0, 0.0}
	si := []float64{0.0, 0.0, 0.0}
	ei := []float64{0.0, 0.0, 0.0}

	el := new(BlurElement)
	el.rows = 5
	el.cols = 5
	el.data = initData
	el.scratch = scratch
	el.northIn = northIn
	el.northOut = northOut
	el.westIn = westIn
	el.westOut = westOut
	el.southIn = southIn
	el.southOut = southOut
	el.eastIn = eastIn
	el.eastOut = eastOut

	// Drain the output buffers
	go func() {
		for i := 0; i < 4; i++ {
			select {
			case <-northOut:
			case <-westOut:
			case <-southOut:
			case <-eastOut:
			}
		}
	}()

	// Make sure there is something waiting to be swapped
	go func() {
		southIn <- si
		northIn <- ni
		eastIn <- ei
		westIn <- wi
	}()

	// Drain outgoing messages
	el.Step()

	expectedData := [][]float64{
		[]float64{0, 0, 0, 0, 0},
		[]float64{0, 0.625, 0.875, 0.625, 0},
		[]float64{0, 0.875, 1.25, 0.875, 0},
		[]float64{0, 0.625, 0.875, 0.625, 0},
		[]float64{0, 0, 0, 0, 0},
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			obs := el.Data()[i][j]
			exp := expectedData[i][j]
			if obs != exp {
				t.Error("Mismatch in blurred data at", i, j, "Expected:", exp, "got", obs)
			}
		}
	}
}

func TestBlurElementSwapsCols(t *testing.T) {
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

	el := new(BlurElement)
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
		northIn <- ni
		westIn <- wi
		southIn <- si
		eastIn <- ei
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

func TestMakeBlurArrayPopulatesData(t *testing.T) {
	globalData := [][]float64{
		[]float64{1, 1, 2, 2},
		[]float64{1, 1, 2, 2},
		[]float64{3, 3, 4, 4},
		[]float64{3, 3, 4, 4},
	}

	widthEls, heightEls := 2, 2
	elRows := len(globalData)/heightEls + 2
	elCols := len(globalData[0])/widthEls + 2
	arr := MakeBlurArray(globalData, widthEls, heightEls)

	elID := 1
	for _, r := range arr.Elements() {
		for _, el := range r {
			if len(el.Data()) != elRows {
				t.Error("Size Mismatch")
			}
			if len(el.Data()[0]) != elCols {
				t.Error("Size Mismatch")
			}

			// Zero-valued boundary cells
			for i := 0; i < elCols; i++ {
				if el.Data()[0][i] != 0 {
					t.Error("Nonzero top boundary detected:", el.Data()[0][i])
				}
				if el.Data()[elRows-1][i] != 0 {
					t.Error("Nonzero bottom boundary detected:", el.Data()[elRows-1][i])
				}
			}
			for i := 0; i < elRows; i++ {
				if el.Data()[i][0] != 0 {
					t.Error("Nonzero left boundary detected:", el.Data()[i][0])
				}
				if el.Data()[elCols-1][i] != 0 {
					t.Error("Nonzero right boundary detected:", el.Data()[i][elCols-1])
				}
			}
			// Inner values preserved:
			for i := 1; i < elRows-1; i++ {
				for j := 1; j < elCols-1; j++ {
					obs := el.Data()[i][j]
					exp := float64(elID)
					if obs != exp {
						t.Error("Mismatched inner data! Got", obs, "wanted", exp)
					}
				}
			}

			elID++
		}
	}
}

func TestMakeBlurArrayReturnsViableArray(t *testing.T) {
	elRows := 5
	elCols := 7
	widthEls := 10
	heightEls := 10

	blurData := make([][]float64, elRows*heightEls)
	for i := 0; i < elRows*heightEls; i++ {
		blurData[i] = make([]float64, elCols*widthEls)
	}

	arr := MakeBlurArray(blurData, widthEls, heightEls)
	arr.Step()
}

func BenchmarkBlur(b *testing.B) {
	elRows := 25
	elCols := 25
	heightEls := 10
	widthEls := 10

	blurData := make([][]float64, elRows*heightEls)
	for i := 0; i < elRows; i++ {
		blurData[i] = make([]float64, elCols*widthEls)
	}

	arr := MakeBlurArray(blurData, widthEls, heightEls)
	for i := 0; i < b.N; i++ {
		arr.Step()
	}
}
