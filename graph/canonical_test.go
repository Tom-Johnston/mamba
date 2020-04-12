package graph_test

import (
	"fmt"
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
)

func TestCanonicalIsomorph(t *testing.T) {
	expectedNumber := []int{1, 1, 2, 4, 11, 34, 156, 1044}
	for n := 1; n < 7; n++ {
		edges := make([]byte, (n*(n-1))/2)
		uniqueGraphs := make(map[string]struct{}, 1044)
		for i := 0; i < (1 << uint((n*(n-1))/2)); i++ {
			for j := 0; j < len(edges); j++ {
				edges[j] = byte((i >> uint(j)) & 1)
			}
			g := graph.NewDense(n, edges)
			uniqueGraphs[fmt.Sprint(graph.Graph6Encode(g.InducedSubgraph(graph.CanonicalIsomorph(g))))] = struct{}{}
		}
		if len(uniqueGraphs) != expectedNumber[n] {
			t.Errorf("Wrong number of graphs on %v vertices. Found: %v Expected: %v", n, len(uniqueGraphs), expectedNumber[n])
			t.Fail()
		}
	}
}

//This will need changing if the cell selection changes.
// func TestEquitableRefinementProcedure(t *testing.T) {
// 	edges := make([]uint8, 36) //See
// 	edges[0] = 1
// 	edges[2] = 1
// 	edges[3] = 1
// 	edges[7] = 1
// 	edges[12] = 1
// 	edges[33] = 1
// 	edges[14] = 1
// 	edges[9] = 1
// 	edges[18] = 1
// 	edges[25] = 1
// 	edges[27] = 1
// 	edges[35] = 1
// 	g := DenseGraph{NumberOfVertices: 9, Edges: edges}
// 	dws := make([]degreeWrapper, g.NumberOfVertices)
// 	for i := 0; i < g.NumberOfVertices; i++ {
// 		dws[i].degrees = make([]int, 1)
// 	}
// 	nbs := make([]int, g.NumberOfVertices)
//
// 	op := OrderedPartition{[]int{0, 2, 6, 8, 1, 3, 5, 7, 4}, []int{1, 3, 4, 1}, []int{}, 0}
// 	r := equitableRefinementProcedure(g, op, 1, dws, nbs)
// 	if !IntsEqual(r.Order, []int{0, 2, 6, 8, 5, 7, 1, 3, 4}) || !IntsEqual(r.BinSizes, []int{1, 2, 1, 2, 2, 1}) {
// 		t.Fail()
// 	}
//
// 	op = OrderedPartition{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8}, []int{9}, []int{}, 0}
// 	r = equitableRefinementProcedure(g, op, 1, dws, nbs)
// 	if !IntsEqual(r.Order, []int{0, 2, 6, 8, 1, 3, 5, 7, 4}) || !IntsEqual(r.BinSizes, []int{4, 4, 1}) {
// 		t.Fail()
// 	}
// }
