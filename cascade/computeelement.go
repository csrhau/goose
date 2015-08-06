package scma

// computeElement represents an elent within a SCMA or other array-based
// computing machine
type computeElement struct {
	rank      int
	data      [][]float64
	clockLine chan bool
	inBus     chan []float64
	outBus    chan []float64
}

// InBus returns the input data channel for this element
func (el *computeElement) InBus() chan<- []float64 {
	return el.inBus
}

// OutBus returns the output data channel for this element
func (el *computeElement) OutBus() <-chan []float64 {
	return el.outBus
}

// Rank returns the rank of the ComputeElement within an array
func (el *computeElement) Rank() int {
	return el.rank
}

func (el *computeElement) communicate() {
	outbound, remainder := el.data[0], el.data[1:]
	go func() { el.outBus <- outbound }()
	inbound := <-el.inBus
	el.data = append(remainder, inbound)
}
