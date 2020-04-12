package graph_test

import (
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
)

func TestInducedSubgraph(t *testing.T) {
	controlG, _ := graph.Graph6Decode("Ks@HOo?PGdCK")
	testGraph, _ := graph.Sparse6Decode(":K`ADOccQXK`IaXcQMb")
	inducedControl := controlG.InducedSubgraph([]int{1, 2, 3, 4, 7, 10, 11})
	inducedTest := testGraph.InducedSubgraph([]int{1, 2, 3, 4, 7, 10, 11})
	if !graph.Equal(inducedControl, inducedTest) {
		t.Log(graph.Graph6Encode(inducedControl), graph.Graph6Encode(inducedTest))
		t.Fail()
	}

	inducedControl = controlG.InducedSubgraph([]int{1, 7, 9, 8, 10, 3, 4, 2})
	inducedTest = testGraph.InducedSubgraph([]int{1, 7, 9, 8, 10, 3, 4, 2})
	if !graph.Equal(inducedControl, inducedTest) {
		t.Log(graph.Graph6Encode(inducedControl), graph.Graph6Encode(inducedTest))
		t.Fail()
	}
}
