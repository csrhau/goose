package goose

import "testing"

func TestElements(t *testing.T) {
	elements := []int{1, 2, 5, 10, 1000}
	for _, els := range elements {
		arr := NewComputeArray(els, els, 1)
		arrSize := arr.Elements()
		if arrSize != els {
			t.Error("Expected Elements()", els, "got", arrSize)
		}
	}
}

func TestEqualDomainDecomposition(t *testing.T) {
	elems := []int{3, 5, 7, 11, 13}
	elWidth := 5
	elHeight := 10
	for _, els := range elems {
		arr := NewComputeArray(els, els*elWidth, elHeight)
		for _, el := range arr.elements {
			if el.FixedDimCells() != elHeight {
				t.Error("Mismatched fixed dimension, expected", elHeight, "got", el.FixedDimCells())
			}
			if el.TransDimCells() != elWidth {
				t.Error("Mismatched trans dimension, expected", elWidth, "got", el.TransDimCells())
			}
		}
	}
}

func TestUnequalDomainDecomposition(t *testing.T) {
	elems := []int{3, 5, 7, 11, 13}
	cells := []int{101, 103, 107, 109, 113, 127}
	elHeight := 10
	for _, els := range elems {
		for _, cells := range cells {
			arr := NewComputeArray(els, cells, elHeight)
			cellCount := 0
			for i, el := range arr.elements {
				var expectedCells int
				// Slightly weird formulation because we back-load the cells
				if (els - i - 1) < cells%els {
					expectedCells = 1 + cells/els
				} else {
					expectedCells = cells / els
				}
				obsCells := el.TransDimCells()
				if obsCells != expectedCells {
					t.Error("Cell count mismatch for element", i, "expected", expectedCells, "got", obsCells)
				}
				cellCount += el.TransDimCells()
			}
			if cellCount != cells {
				t.Error("Conservation of cells violated! Expected", cells, "got", cellCount)
			}
		}
	}
}

func TestAddingData(t *testing.T) {
}

func TestDrainingData(t *testing.T) {
}

func TestDataThroughput(t *testing.T) {

}
