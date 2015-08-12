package goose

// ComputeElement describes elements in a computational array
type ComputeElement interface {
	Data() [][]float64
	Step()
}
