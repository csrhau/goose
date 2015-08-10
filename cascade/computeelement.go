package cascade

// ComputeElement represents an elent within a cascade or other array-based
// computing machine
type ComputeElement struct {
	data      [][]float64
	clockLine chan bool
	inBus     chan []float64
	outBus    chan []float64
}

// FixedDimCells returns the number of cells along the fixed axis
// orthogonal to the translating axis. This corresponds to the length
// of the slices which will be transmitted.
func (el *ComputeElement) FixedDimCells() int {
	return len(el.data[0])
}

// TransDimCells returns the number of cells along the translating axis
func (el *ComputeElement) TransDimCells() int {
	return len(el.data)
}

// InBus returns the input data channel for this element
func (el *ComputeElement) InBus() chan<- []float64 {
	return el.inBus
}

// OutBus returns the output data channel for this element
func (el *ComputeElement) OutBus() <-chan []float64 {
	return el.outBus
}

func (el *ComputeElement) communicate() {
	outbound, remainder := el.data[0], el.data[1:]
	go func() { el.outBus <- outbound }()
	inbound := <-el.inBus
	el.data = append(remainder, inbound)
}

func NewComputeElement(x, y int, cl chan bool, ib, ob chan []float64) ComputeElement {
	data := make([][]float64, x)
	for i := 0; i < x; i++ {
		data[i] = make([]float64, y)
	}
	return ComputeElement{data, cl, ib, ob}
}
