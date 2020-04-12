package graph_test

import (
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/graph/search"
	"github.com/Tom-Johnston/mamba/ints"
)

func TestGraph6(t *testing.T) {
	for i := 1; i <= 8; i++ {
		output := make(chan *graph.DenseGraph)
		go search.All(i, output, 0, 1)
		for g := range output {
			s := graph.Graph6Encode(g)
			h, err := graph.Graph6Decode(s)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			if !graph.Equal(g, h) {
				t.Logf("Graphs not equal")
				t.FailNow()
			}
		}
	}
}

func TestSparse6Encode(t *testing.T) {
	testGraphs := [][]string{
		[]string{"Ks@HOo?PGdCK", ":K`ADOccQXK`IaXcQMb"},
		[]string{"OsaBA`GP@`dIHWEcas_]O", ":O`ACGPDC[QPJGYCqG\\KafPK`ckeSqDsIWyn"},
		[]string{"J?AKagjXfo?", ":Ji?c@pEUPBFaGhg@CKf"},
	}
	for _, s := range testGraphs {
		g, err := graph.Graph6Decode(s[0])
		if err != nil {
			t.Error("Failed to decode Graph6 format.")
		}
		if c := graph.Sparse6Encode(g); s[1] != c {
			t.Log(graph.Sparse6Decode(c))
			t.Logf("Expected: %v Found: %v\n", s[1], c)
			t.Fail()
		}
	}
}

func TestSparse6Decode(t *testing.T) {
	testGraphs := [][]string{
		[]string{"Ks@HOo?PGdCK", ":K`ADOccQXK`IaXcQMb"},
		[]string{"OsaBA`GP@`dIHWEcas_]O", ":O`ACGPDC[QPJGYCqG\\KafPK`ckeSqDsIWyn"},
		[]string{"J?AKagjXfo?", ":Ji?c@pEUPBFaGhg@CKf"},
	}
	for _, s := range testGraphs {
		g, err := graph.Sparse6Decode(s[1])
		if err != nil {
			t.Error("Failed to decode Sparse6 format.")
		}
		if c := graph.Graph6Encode(g); s[0] != c {
			t.Logf("Expected: %v Found: %v\n", s[0], c)
			t.Fail()
		}
	}
}

func TestPruferEncode(t *testing.T) {
	g := graph.NewDense(6, []byte{0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1})
	code := []int{3, 3, 3, 4}
	if c := graph.PruferEncode(g); !ints.Equal(c, code) {
		t.Logf("Found: %v Expected: %v", c, code)
		t.Fail()
	}
}

func TestPruferDecode(t *testing.T) {
	g := graph.NewDense(6, []byte{0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1})
	code := []int{3, 3, 3, 4}
	if f := graph.PruferDecode(code); !graph.Equal(f, g) {
		t.Logf("Found: %v Expected: %v", f, g)
		t.Fail()
	}
}
