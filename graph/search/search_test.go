package search

import (
	"runtime"
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/ints"
)

func TestAll(t *testing.T) {
	truthData := []int{1, 1, 2, 4, 11, 34, 156, 1044, 12346, 274668, 12005168, 1018997864, 165091172592, 50502031367952}
	maxSize := 9 //Largest graphs we will search for.
	numberFound := make([]int, maxSize+1)
	for size := 0; size <= maxSize; size++ {
		output := make(chan *graph.DenseGraph, 1)
		go All(size, output, 0, 1)

		for range output {
			numberFound[size]++
		}
	}
	t.Log(numberFound)
	if !ints.Equal(numberFound, truthData[:maxSize+1]) {
		t.Fail()
	}
}

func TestAllParallel(t *testing.T) {
	truthData := []int{1, 1, 2, 4, 11, 34, 156, 1044, 12346, 274668, 12005168, 1018997864, 165091172592, 50502031367952}
	maxSize := 9 //Largest graphs we will search for.
	numberFound := make([]int, maxSize+1)
	for size := 0; size <= maxSize; size++ {
		output := make(chan *graph.DenseGraph, 1)
		go AllParallel(size, output, runtime.GOMAXPROCS(0))

		for range output {
			numberFound[size]++
		}
	}
	t.Log(numberFound)
	if !ints.Equal(numberFound, truthData[:maxSize+1]) {
		t.Fail()
	}
}

func BenchmarkAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output := make(chan *graph.DenseGraph, 1)
		go All(10, output, 0, 1)
		for range output {
		}
	}
}
