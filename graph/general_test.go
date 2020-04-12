package graph_test

import (
	"testing"

	"github.com/Tom-Johnston/mamba/graph/search"
	"github.com/Tom-Johnston/mamba/ints"
	"github.com/Tom-Johnston/mamba/sortints"

	"github.com/Tom-Johnston/mamba/graph"
)

func TestDegeneracy(t *testing.T) {
	//truthData[n][k] contais the number of graphs of size n + 1 with chromatic numer k
	truthData := make([][]int, 10)
	truthData[0] = []int{1}
	truthData[1] = []int{1}
	truthData[2] = []int{1, 1}
	truthData[3] = []int{1, 2, 1}
	truthData[4] = []int{1, 5, 4, 1}
	truthData[5] = []int{1, 9, 18, 5, 1}
	truthData[6] = []int{1, 19, 85, 43, 7, 1}
	truthData[7] = []int{1, 36, 471, 442, 85, 8, 1}
	truthData[8] = []int{1, 75, 3378, 6979, 1758, 144, 10, 1}
	truthData[9] = []int{1, 152, 31782, 166258, 70811, 5421, 231, 11, 1}
	for i := 1; i <= 8; i++ {
		output := make(chan *graph.DenseGraph)
		foundData := make([]int, i)
		go search.All(i, output, 0, 1)
		for g := range output {
			d, order := graph.Degeneracy(g)
			for i, v := range order {
				orderD := 0
				neighbours := g.Neighbours(v)
				for j := 0; j < i; j++ {
					if sortints.ContainsSingle(neighbours, order[j]) {
						orderD++
					}
				}
				if orderD > d {
					t.Fail()
				}
			}
			foundData[d]++
		}
		if !ints.Equal(foundData, truthData[i]) {
			t.Log(foundData)
			t.Log(truthData)
			t.Fail()
		}
	}
}

func TestBiconnectedComponents(t *testing.T) {
	maxSize := 8
	truthData := []int{0, 0, 1, 1, 3, 10, 56, 468, 7123, 194066, 9743542, 900969091}
	foundData := make([]int, maxSize+1)
	for i := 0; i <= maxSize; i++ {
		output := make(chan *graph.DenseGraph)
		go search.All(i, output, 0, 1)
		for g := range output {
			c, _ := graph.BiconnectedComponents(g)
			if len(c) == 1 && len(c[0]) > 1 {
				foundData[g.N()]++
			}
		}
	}
	if !ints.Equal(foundData, truthData[:maxSize+1]) {
		t.Log(foundData)
		t.Log(truthData[:maxSize+1])
		t.Fail()
	}
}
