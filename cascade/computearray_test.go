package scma

import "testing"

func TestElements(t *testing.T) {
	sizes := []int{1, 2, 5, 10, 1000}
	for _, size := range sizes {
		scma := NewComputeArray(size)
		scmaSize := scma.Elements()
		if scmaSize != size {
			t.Error("Expected Elements() ", size, ", got ", scmaSize)
		}
	}
}

func TestContiguousElements(t *testing.T) {
	sizes := []int{1, 2, 5, 10, 1000}
	for _, size := range sizes {
		scma := NewComputeArray(size)
		for i := 0; i < size; i++ {
			rank := scma.elements[i].Rank()
			if rank != i {
				t.Error("Unsequential Rank Detected! Expected ", i, " got ", rank)
			}
		}
	}
}

func TestAddingData(t *testing.T) {

}

func TestDrainingData(t *testing.T) {

}
