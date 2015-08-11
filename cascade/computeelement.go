package goose

// ComputeElement describes elements in a computational array
type ComputeElement interface {
	ClockLine() chan<- bool
	Data() [][]float64
	Step()
}
