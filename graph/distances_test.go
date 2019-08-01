package graph_test

import (
	"testing"

	"github.com/Tom-Johnston/gigraph/graph"
	"github.com/Tom-Johnston/gigraph/graph/search"
)

//The truth data was generated using the function but while there were other tests.
func TestGirth(t *testing.T) {
	truthData := make([][]int, 10)
	truthData[0] = []int{0}
	truthData[1] = []int{0, 0}
	truthData[2] = []int{0, 0, 0}
	truthData[3] = []int{0, 0, 0, 1}
	truthData[4] = []int{0, 0, 0, 4, 1}
	truthData[5] = []int{0, 0, 0, 20, 3, 1}
	truthData[6] = []int{0, 0, 0, 118, 15, 2, 1}
	truthData[7] = []int{0, 0, 0, 937, 59, 8, 2, 1}
	truthData[8] = []int{0, 0, 0, 11936, 296, 26, 9, 2, 1}
	for i := 0; i <= 8; i++ {
		output := make(chan *graph.DenseGraph)
		foundData := make([]int, i+1)

		go search.All(i, output, 0, 1)
		for g := range output {
			girth := graph.Girth(g)
			if girth > 0 {
				foundData[girth]++
			}
		}
		if !graph.IntsEqual(foundData, truthData[i]) {
			t.Log(foundData)
			t.Log(truthData)
			t.Fail()
		}
	}
	//Test a few known graphs.
	g := graph.CompleteGraph(7)
	if graph.Girth(g) != 3 {
		t.Errorf("Graph: CompleteGraph(7) Found: %v Expected: 3", graph.Girth(g))
	}

	g = graph.NewDense(5, nil)
	if graph.Girth(g) != -1 {
		t.Errorf("Graph: EmptyGraph(5) Found: %v Expected: -1", graph.Girth(g))
	}

	g = graph.Cycle(6)
	if graph.Girth(g) != 6 {
		t.Errorf("Graph: Cycle(6) Found: %v Expected: 6", graph.Girth(g))
	}
}

func TestDiameter(t *testing.T) {
	truthData := make([][]int, 10)
	truthData[0] = []int{1}
	truthData[1] = []int{1}
	truthData[2] = []int{0, 1}
	truthData[3] = []int{0, 1, 1}
	truthData[4] = []int{0, 1, 4, 1}
	truthData[5] = []int{0, 1, 14, 5, 1}
	truthData[6] = []int{0, 1, 59, 43, 8, 1}
	truthData[7] = []int{0, 1, 373, 387, 82, 9, 1}
	truthData[8] = []int{0, 1, 4154, 5797, 1027, 125, 12, 1}
	for i := 1; i <= 7; i++ {
		output := make(chan *graph.DenseGraph)
		foundData := make([]int, i)

		go search.All(i, output, 0, 1)
		for g := range output {
			diam := graph.Diameter(g)
			if diam > -1 {
				foundData[diam]++
			}
		}
		if !graph.IntsEqual(foundData, truthData[i]) {
			t.Log(foundData)
			t.Log(truthData)
			t.Fail()
		}
	}
	//Test a few known graphs.
	g := graph.NewDense(0, nil)
	if graph.Diameter(g) != 0 {
		t.Errorf("Graph: Empty(0) Found: %v Expected: 1", graph.Diameter(g))
	}

	g = graph.CompleteGraph(7)
	if graph.Diameter(g) != 1 {
		t.Errorf("Graph: CompleteGraph(7) Found: %v Expected: 1", graph.Diameter(g))
	}

	g = graph.NewDense(5, nil)
	if graph.Diameter(g) != -1 {
		t.Errorf("Graph: EmptyGraph(5) Found: %v Expected: -1", graph.Diameter(g))
	}

	g = graph.Cycle(6)
	if graph.Diameter(g) != 3 {
		t.Errorf("Graph: Cycle(6) Found: %v Expected: 3", graph.Diameter(g))
	}

	g = graph.Star(5)
	if graph.Diameter(g) != 2 {
		t.Errorf("Graph: Star(5) Found: %v Expected: 2", graph.Diameter(g))
	}
}

func TestRadius(t *testing.T) {
	truthData := make([][]int, 10)
	truthData[0] = []int{1}
	truthData[1] = []int{1}
	truthData[2] = []int{0, 1}
	truthData[3] = []int{0, 2}
	truthData[4] = []int{0, 4, 2}
	truthData[5] = []int{0, 11, 10}
	truthData[6] = []int{0, 34, 76, 2}
	truthData[7] = []int{0, 156, 682, 15}
	truthData[8] = []int{0, 1044, 9864, 207, 2}
	for i := 0; i <= 8; i++ {
		output := make(chan *graph.DenseGraph)
		foundData := make([]int, (i/2)+1)

		go search.All(i, output, 0, 1)
		for g := range output {
			diam := graph.Radius(g)
			if diam > -1 {
				foundData[diam]++
			}
		}
		if !graph.IntsEqual(foundData, truthData[i]) {
			t.Log(foundData)
			t.Log(truthData)
			t.Fail()
		}
	}
	//Test a few known graphs.
	g := graph.CompleteGraph(7)
	if graph.Radius(g) != 1 {
		t.Errorf("Graph: CompleteGraph(7) Found: %v Expected: 1", graph.Radius(g))
	}

	g = graph.NewDense(5, nil)
	if graph.Radius(g) != -1 {
		t.Errorf("Graph: EmptyGraph(5) Found: %v Expected: -1", graph.Radius(g))
	}

	g = graph.Cycle(6)
	if graph.Radius(g) != 3 {
		t.Errorf("Graph: Cycle(6) Found: %v Expected: 3", graph.Radius(g))
	}

	g = graph.Star(5)
	if graph.Radius(g) != 1 {
		t.Errorf("Graph: Star(5) Found: %v Expected: 1", graph.Radius(g))
	}
}
