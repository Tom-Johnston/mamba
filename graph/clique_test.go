package graph_test

import (
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/graph/search"
	"github.com/Tom-Johnston/mamba/ints"
)

func TestCliqueNumber(t *testing.T) {
	//Check the number of graphs on vertices with each clique number. This includes the empty graph.
	truthData := make([][]int, 11)
	truthData[0] = []int{1}
	truthData[1] = []int{0, 1}
	truthData[2] = []int{0, 1, 1}
	truthData[3] = []int{0, 1, 2, 1}
	truthData[4] = []int{0, 1, 6, 3, 1}
	truthData[5] = []int{0, 1, 13, 15, 4, 1}
	truthData[6] = []int{0, 1, 37, 82, 30, 5, 1}
	truthData[7] = []int{0, 1, 106, 578, 301, 51, 6, 1}
	truthData[8] = []int{0, 1, 409, 6021, 4985, 842, 80, 7, 1}
	truthData[9] = []int{0, 1, 1896, 101267, 142276, 27107, 1995, 117, 8, 1}
	truthData[10] = []int{0, 1, 12171, 2882460, 7269487, 1724440, 112225, 4210, 164, 9, 1}
	for i := 0; i <= 8; i++ {
		output := make(chan *graph.DenseGraph)
		foundData := make([]int, i+1)
		go search.All(i, output, 0, 1)
		for g := range output {
			foundData[graph.CliqueNumber(g)]++
		}
		if !ints.Equal(foundData, truthData[i]) {
			t.Log(foundData)
			t.Log(truthData)
			t.Fail()
		}
	}
	//Check the value for a couple of fixed graphs.
	if c := graph.CliqueNumber(graph.CompleteGraph(5)); c != 5 {
		t.Logf("Complete Graph Expected: 5 Found: %v\n", c)
		t.Fail()
	}

	if c := graph.CliqueNumber(graph.Cycle(5)); c != 2 {
		t.Logf("Complete Graph Expected: 2 Found: %v\n", c)
		t.Fail()
	}

}

func TestIndependenceNumber(t *testing.T) {
	testGraphs := make(map[string]int)
	testGraphs["Ks@HOo?PGdCK"] = 5
	testGraphs["OsaBA`GP@`dIHWEcas_]O"] = 5
	testGraphs["J?AKagjXfo?"] = 5
	testGraphs["IsP@OkWHG"] = 4

	for g6, trueNumber := range testGraphs {
		g, err := graph.Graph6Decode(g6)
		if err != nil {
			t.Log(err)
			t.Fail()
			continue
		}
		testNumber := graph.IndependenceNumber(g)
		if testNumber != trueNumber {
			t.Fail()
		}
	}

}
