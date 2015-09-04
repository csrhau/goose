package goose

import "testing"

func TestArrayIterations(t *testing.T) {
	els := []ComputeElement{
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
	}
	ar := NewComputeArray(els, 2, 2)
	for i := 1; i < 10; i++ {
		ar.Step()
		for _, el := range ar.Elements() {
			for _, r := range el.Data() {
				for _, v := range r {
					if v != float64(i) {
						t.Error("Counting element data error")
					}
				}
			}
		}
	}
}

func TestArrayClocking(t *testing.T) {
	els := []ComputeElement{
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
		new(CountingElement),
	}
	ar := NewComputeArray(els, 2, 2)
	clk := make(chan bool)
	go ar.Run(clk)
	for i := 1; i < 10; i++ {
		clk <- true
		for _, el := range ar.Elements() {
			for _, r := range el.Data() {
				for _, v := range r {
					if v != float64(i) {
						t.Error("Counting element data error")
					}
				}
			}
		}
	}
	close(clk)
}

func TestArrayLayout(t *testing.T) {
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			els := make([]ComputeElement, i*j)
			for ij := 0; ij < i*j; ij++ {
				els[ij] = new(CountingElement)
			}
			arr := NewComputeArray(els, i, j)
			rows, cols := arr.Layout()
			if rows != i {
				t.Error("Row array mismatch, expected", i, "got", rows)
			}
			if cols != j {
				t.Error("Row array mismatch, expected", j, "got", cols)
			}
		}
	}
}

func TestArrayElementAt(t *testing.T) {
	elsVertical := 8
	elsHorizontal := 13
	arr := MakeCountingArray(elsVertical, elsHorizontal, 1, 1)
	idx := 0
	for row := 0; row < elsVertical; row++ {
		for col := 0; col < elsHorizontal; col, idx = col+1, idx+1 {
			if arr.ElementAt(row, col).Data()[0][0] != float64(idx) {
				t.Error("Problem in ElementAt")
			}
		}
	}
}

func TestArrayShape(t *testing.T) {
	for elsVertical := 1; elsVertical < 10; elsVertical++ {
		for elsHorizontal := 1; elsHorizontal < 10; elsHorizontal++ {
			for elRows := 1; elRows < 10; elRows++ {
				for elCols := 1; elCols < 10; elCols++ {
					arr := MakeCountingArray(elsVertical, elsHorizontal, elRows, elCols)
					gr, gc := arr.Shape()
					if gr != elsVertical*elRows {
						t.Error("Error in array shape, rows")
					}
					if gc != elsHorizontal*elCols {
						t.Error("Error in array shape, cols")
					}
				}
			}
		}
	}
}

func TestArrayDataSimple(t *testing.T) {
	arr := MakeCountingArray(3, 5, 1, 1)
	exp := [][]float64{
		[]float64{0, 1, 2, 3, 4},
		[]float64{5, 6, 7, 8, 9},
		[]float64{10, 11, 12, 13, 14},
	}
	act := arr.Data()
	for i, r := range exp {
		for j, c := range r {
			if c != act[i][j] {
				t.Error("DataMismatch!")
			}
		}
	}
}

/* TODO
func TestUnevenConstructionPanics(t *testing.T) {
}


func TestArrayData(t *testing.T) {
}
*/
