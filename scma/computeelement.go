package scma

// computeElement represents an element within a SCMA or other array-based
// computing machine

type computeElement struct {
	rank      int
	data      [][]float64
	clockLine chan bool
	// TODO I think these should be North, South.
	dataLineEast chan []float64 // In
	dataLineWest chan []float64 // Out
}

// Rank returns the rank of the ComputeElement within an array
func (elem *computeElement) Rank() int {
	return elem.rank
}

func (elem *computeElement) communicate() {
	var outbound, inbound []float64
	inbound = <-elem.dataLineEast
	outbound, elem.data = elem.data[0], append(elem.data[1:], inbound)
	elem.dataLineWest <- outbound
}
