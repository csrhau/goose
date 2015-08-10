package goose

// ComputeElement describes elements in a computational array
type ComputeElement interface {
	ClockLine() chan<- bool
	State() [][]float64
	Step()
}
