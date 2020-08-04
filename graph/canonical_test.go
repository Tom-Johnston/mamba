package graph_test

import (
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
)

func TestCanonicalIsomorph(t *testing.T) {
	expectedNumber := []int{1, 1, 2, 4, 11, 34, 156, 1044, 12346}
	for n := 0; n < 7; n++ {
		edges := make([]byte, (n*(n-1))/2)
		uniqueGraphs := make(map[string]struct{}, 1044)
		for i := 0; i < (1 << uint((n*(n-1))/2)); i++ {
			for j := 0; j < len(edges); j++ {
				edges[j] = byte((i >> uint(j)) & 1)
			}
			g := graph.NewDense(n, edges)
			uniqueGraphs[graph.Graph6Encode(g.InducedSubgraph(graph.CanonicalIsomorph(g)))] = struct{}{}
		}
		if len(uniqueGraphs) != expectedNumber[n] {
			t.Errorf("Wrong number of graphs on %v vertices. Found: %v Expected: %v", n, len(uniqueGraphs), expectedNumber[n])
			t.Fail()
		}
	}
}

func BenchmarkCanonicalIsomorph(b *testing.B) {
	graphs := []string{"S}GOOOE@?C?K?O?E_A??S?@C??_?@G??[",
		"S{S_gOD?_A?E?E?B??O?A??G??_??w??w",
		"S{S__OE@?C?H?O?G_A??`?AO?C??DG??w",
		"S{S__OC@?C?O?O?O?D_@??EC?OO?M??EG",
		"S{O_o_H@?G?O?S?K?@??P?@_?C_?A_??[",
		"S{O_o_G@?G?Q?Y?O?C??_?AO?D??@o?@g",
		"S{O___IA?K?_?W?U?AO?G?@???_??w??w",
		"S{O___IA?G?_?a?P?G?@A?E??O_?Gg?DO"}
	for i := 0; i < b.N; i++ {
		for _, g6 := range graphs {
			g, _ := graph.Graph6Decode(g6)
			graph.CanonicalIsomorph(g)
		}
	}
}
