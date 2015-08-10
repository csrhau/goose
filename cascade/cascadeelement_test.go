package goose

import "testing"

func TestCommunicationPassthrough(t *testing.T) {
	data := [][]float64{
		make([]float64, 3),
		make([]float64, 3),
		make([]float64, 3),
	}
	upstream := make(chan []float64)
	downstream := make(chan []float64)
	elem := NewCascadeElement(3, 3, nil, upstream, downstream)

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
