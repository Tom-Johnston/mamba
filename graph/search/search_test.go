package search

import (
	"testing"

	. "github.com/Tom-Johnston/gigraph/graph"
)

func TestAll(t *testing.T) {
	truthData := []int{1, 1, 2, 4, 11, 34, 156, 1044, 12346, 274668, 12005168, 1018997864, 165091172592, 50502031367952}
	maxSize := 9 //Largest graphs we will search for.
	numberFound := make([]int, maxSize+1)
	for size := 0; size <= maxSize; size++ {
		output := make(chan *DenseGraph, 1)
		go All(size, output, 0, 1)

		for _ = range output {
			numberFound[size]++
		}
	}
	t.Log(numberFound)
	if !IntsEqual(numberFound, truthData[:maxSize+1]) {
		t.Fail()
	}
}

func TestAllParallel(t *testing.T) {
	truthData := []int{1, 1, 2, 4, 11, 34, 156, 1044, 12346, 274668, 12005168, 1018997864, 165091172592, 50502031367952}
	maxSize := 9 //Largest graphs we will search for.
	numberFound := make([]int, maxSize+1)
	for size := 0; size <= maxSize; size++ {
		output := make(chan *DenseGraph, 1)
		go AllParallel(size, output)

		for _ = range output {
			numberFound[size]++
		}
	}
	t.Log(numberFound)
	if !IntsEqual(numberFound, truthData[:maxSize+1]) {
		t.Fail()
	}
}
