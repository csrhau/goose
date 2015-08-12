package goose

import "fmt"

// RowSwapElement is a dummy ComputeElement which exchanges its top and bottom row
type RowSwapElement struct {
	data              [][]float64
	northIn, northOut chan []float64
	southIn, southOut chan []float64
}

func (el *RowSwapElement) Data() [][]float64 {
	return el.data
}

func (el *RowSwapElement) Swap() {
	sent := make(chan struct{})
	// Send - TODO benchmark with select
	fmt.Println("Swap was called!")
	go func() {
		fmt.Println("Send was triggered")
		el.northOut <- el.data[0]
		fmt.Println("Sent.. something")
		el.southOut <- el.data[len(el.data)-1]
		fmt.Println("Sent.. something")
		close(sent)
	}()

	var northRecv, southRecv []float64
	// Receive - TODO benchmark without select, just reversed ordering
	for i := 0; i < 2; i++ {
		fmt.Println("Await data..")
		select {
		case northRecv = <-el.northIn:
			fmt.Println("GOT DATA FROM THE north!")
		case southRecv = <-el.southIn:
			fmt.Println("GOT DATA FROM THE south!")
		}
	}

	// Update
	<-sent
	el.data[0] = northRecv
	el.data[len(el.data)-1] = southRecv
}

func (el *RowSwapElement) Step() {
	el.Swap()
}
