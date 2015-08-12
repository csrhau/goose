package goose

import "fmt"

// ComputeElement describes elements in a computational array
type ComputeElement interface {
	Data() [][]float64
	Step()
}

func greet(c ComputeElement) {
	fmt.Println(c)
}
