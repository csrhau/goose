package goose

// CartesianSwapElement is a dummy ComputeElement which performs a standard
// von neuman neighbourhood boundary exchange
type CartesianSwapElement struct {
	clockLine     chan bool
	dataLineNorth chan []float64
	dataLineEast  chan []float64
	dataLineSouth chan []float64
	dataLineWest  chan []float64
	data          [][]float64
}

func (el *CartesianSwapElement) ClockLine() chan<- bool {
	return el.clockLine
}

func Data(el *CartesianSwapElement) [][]float64 {
	return el.data
}

func (el *CartesianSwapElement) Communicate() {
	go func() {
		el.dataLineNorth <- []float64{1, 2, 3, 4}
		el.dataLineEast <- []float64{1, 2, 3, 4}
		el.dataLineSouth <- []float64{1, 2, 3, 4}
		el.dataLineWest <- []float64{1, 2, 3, 4}
	}()
	<-el.dataLineNorth
	<-el.dataLineEast
	<-el.dataLineSouth
	<-el.dataLineWest
}

func (el *CartesianSwapElement) Step() {
	for range el.clockLine {
		el.Communicate()
	}
}
