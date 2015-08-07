package cascade

import "testing"

func TestRank(t *testing.T) {
	ranks := []int{1, 2, 5, 10, 1000}
	for _, rank := range ranks {
		elem := computeElement{rank, nil, nil, nil, nil}
		elemRank := elem.Rank()
		if elemRank != rank {
			t.Error("Expected Rank()", rank, "got", elemRank)
		}
	}
}

func TestCommunicatePresevesInternalData(t *testing.T) {
	data := [][]float64{
		[]float64{1.1, 2.2, 3.3},
		[]float64{4.4, 5.5, 6.6},
		[]float64{7.7, 8.8, 9.9},
	}
	upstream := make(chan []float64)
	downstream := make(chan []float64)
	inSlice := make([]float64, 3)
	elem := computeElement{0, data, nil, upstream, downstream}
	for i := 0; i < 3; i++ {
		go func() { elem.InBus() <- inSlice }()
		elem.communicate()
		observed := <-elem.OutBus()
		expected := data[i]
		lo, le := len(observed), len(expected)
		if lo != le {
			t.Error("Length of expected", lo, "and received", le, "differ!")
		} else {
			for i, e := range expected {
				if observed[i] != e { // intentionally checking bitwise equality
					t.Error("Communications Mismatch! Expected", e, "got", observed[i])
				}
			}
		}
	}
}

func TestCommunicationPassthrough(t *testing.T) {
	data := [][]float64{
		make([]float64, 3),
		make([]float64, 3),
		make([]float64, 3),
	}
	upstream := make(chan []float64)
	downstream := make(chan []float64)
	elem := computeElement{0, data, nil, upstream, downstream}

	payload := [][]float64{
		[]float64{1.0, 2.0, 3.0},
		[]float64{4.0, 5.0, 6.0},
		[]float64{7.0, 8.0, 9.0},
	}

	// Queue up test payload, followed by initial dummy data
	go func() {
		for _, v := range payload {
			elem.InBus() <- v
		}
		for _, v := range data {
			elem.InBus() <- v
		}
	}()

	// Drain initial dummy data
	for i := 0; i < len(data); i++ {
		elem.communicate()
		<-elem.OutBus()
	}

	// Drain payload
	for i := 0; i < len(payload); i++ {
		elem.communicate()
		rcv := <-elem.OutBus()
		for j, v := range rcv {
			e := payload[i][j]
			if e != v {
				t.Error("Communications Mismatch! Expected", e, "got", v)
			}
		}
	}
}
