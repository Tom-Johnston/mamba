package graph_test

import (
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/graph/search"
	"github.com/Tom-Johnston/mamba/ints"
)

func TestPlanarGraph(t *testing.T) {
	maxSize := 8
	truthData := []int{1, 1, 2, 4, 11, 33, 142, 822, 6966, 79853, 1140916}
	foundData := make([]int, maxSize+1)
	// seen := make(map[string]graph.DenseGraph)
	for i := 0; i <= maxSize; i++ {
		output := make(chan *graph.DenseGraph)
		go search.All(i, output, 0, 1)
		for g := range output {
			if graph.IsPlanar(g) {
				foundData[g.N()]++
			}
		}
	}
	if !ints.Equal(foundData, truthData[:maxSize+1]) {
		t.Log(foundData)
		t.Log(truthData[:maxSize+1])
		t.Fail()
	}
	if !graph.IsPlanar(graph.CompleteGraph(4)) {
		t.Log("K4 - Found: false Expected: true")
		t.Fail()
	}
	if graph.IsPlanar(graph.CompleteGraph(5)) {
		t.Log("K5 - Found: true Expected: false")
		t.Fail()
	}
	if graph.IsPlanar(graph.CompletePartiteGraph(3, 3)) {
		t.Log("K_{3,3} - Found: true Expected: false")
	}
}
