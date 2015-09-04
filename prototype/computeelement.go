package goose

// ComputeElement describes elements in a computational array
type ComputeElement interface {
	Data() [][]float64
	Shape() (int, int)
	Step()
}

func (el *ComputeElement) Shape() (int, int) {
	return len(el.Data()), len(el.Data()[0])
}
