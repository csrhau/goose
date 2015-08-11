package goose

// CountingElement is a dummy ComputeElement which simply counts cycles
type CountingElement struct {
	cycles    int
	clockLine chan bool
	data      [][]float64
}

func (el *CountingElement) Cycles() int {
	return el.cycles
}

func (el *CountingElement) Data() [][]float64 {
	return el.data
}

func (el *CountingElement) ClockLine() chan<- bool {
	return el.clockLine
}

func (el *CountingElement) Step() {
	for range el.clockLine {
		el.cycles++
	}
}
