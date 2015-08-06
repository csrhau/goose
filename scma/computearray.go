package scma

// ComputeArray represents a SCMA or other array-based computing machine
type ComputeArray struct {
	elements   []computeElement
	clockLines []chan bool
	commsLines []chan []float64
}

// Elements returns the number of ComputeElements in this array
func (arr *ComputeArray) Elements() int {
	return len(arr.elements)
}

func (arr *ComputeArray) Tick() {
	for _, line := range arr.clockLines {
		line <- true
	}
}

// NewComputeArray builds a ComputeArray of the given dimension
func NewComputeArray(elements int) *ComputeArray {
	arr := new(ComputeArray)
	arr.elements = make([]computeElement, elements)
	arr.clockLines = make([]chan bool, elements)
	arr.commsLines = make([]chan []float64, elements)

	for i := 0; i < elements; i++ {
		arr.clockLines[i] = make(chan bool)
		arr.commsLines[i] = make(chan []float64)
	}

	for i := 0; i < elements; i++ {
		element := computeElement{rank: i,
			data:         nil,
			clockLine:    arr.clockLines[i],
			dataLineEast: arr.commsLines[i],
			dataLineWest: arr.commsLines[(i+1)%elements]}
		arr.elements[i] = element
	}
	return arr
}
