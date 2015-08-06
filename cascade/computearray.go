package scma

// ComputeArray represents a SCMA or other array-based computing machine

// Phase represents the computational state of a ComputeElement
type Phase int

// These constants represent the possible states of a ComputeElement
const (
	Halt Phase = iota
	Fill
	Process
	Drain
)

// ComputeArray models an SCMA or other array based computing machine
type ComputeArray struct {
	phase      Phase
	elements   []computeElement
	inlet      chan []float64
	outlet     chan []float64
	clockLines []chan bool
	commsLines []chan []float64
}

// Phase returns the current state of the array, e.g. halted, processing etc
func (arr *ComputeArray) Phase() Phase {
	return arr.phase
}

// Elements returns the number of ComputeElements in this array
func (arr *ComputeArray) Elements() int {
	return len(arr.elements)
}

// Tick sends a clock pulse to all computeElements in this array
func (arr *ComputeArray) Tick() {
	for _, line := range arr.clockLines {
		line <- true
	}
}

// Inlet is the channel which receives data from upstream
func (arr *ComputeArray) Inlet() chan<- []float64 {
	return arr.inlet
}

// Outlet is the channel which passes data downstream
func (arr *ComputeArray) Outlet() <-chan []float64 {
	return arr.outlet
}

// NewComputeArray builds a cyclic ComputeArray of the given dimension
func NewComputeArray(elements int) *ComputeArray {
	arr := new(ComputeArray)
	arr.phase = Halt
	arr.elements = make([]computeElement, elements)
	arr.inlet = make(chan []float64)
	arr.outlet = make(chan []float64)
	arr.clockLines = make([]chan bool, elements)
	arr.commsLines = make([]chan []float64, elements+2)
	arr.commsLines[0] = arr.inlet
	arr.commsLines[elements+1] = arr.outlet

	for i := 0; i < elements; i++ {
		arr.clockLines[i] = make(chan bool)
		arr.commsLines[i+1] = make(chan []float64)
	}

	arr.inlet = arr.commsLines[0]
	arr.outlet = arr.commsLines[0]

	for i := 0; i < elements; i++ {
		element := computeElement{
			rank:      i,
			data:      nil,
			clockLine: arr.clockLines[i],
			inBus:     arr.commsLines[i],
			outBus:    arr.commsLines[(i+1)%elements],
		}
		arr.elements[i] = element
	}
	return arr
}
