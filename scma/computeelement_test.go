package scma

import "testing"

func TestRank(t *testing.T) {
  ranks := []int {1, 2, 5, 10, 1000}
  for _, rank := range ranks {
    elem := computeElement{rank, nil, nil, nil, nil}
    elemRank := elem.Rank()
    if elemRank != rank {
      t.Error("Expected Rank() ", rank, " got ", elemRank)

    }
  }
}


func TestCommunicateSingal(t *testing.T) {
  t.Error("Not Yet Implemented!")
}



