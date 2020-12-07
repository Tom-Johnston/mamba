package search

import (
	"bytes"
	"testing"

	"github.com/Tom-Johnston/mamba/graph"
	"github.com/Tom-Johnston/mamba/ints"
)

func TestAll(t *testing.T) {
	truthData := []int{1, 1, 2, 4, 11, 34, 156, 1044, 12346, 274668, 12005168, 1018997864, 165091172592, 50502031367952}
	maxSize := 9 //Largest graphs we will search for.
	numberFound := make([]int, maxSize+1)
	for n := 0; n <= maxSize; n++ {
		iter := All(n, 0, 1)
		for iter.Next() {
			numberFound[n]++
		}
	}
	t.Log(numberFound)
	if !ints.Equal(numberFound, truthData[:maxSize+1]) {
		t.Fail()
	}
}

func TestAllConcurrent(t *testing.T) {
	truthData := []int{1, 1, 2, 4, 11, 34, 156, 1044, 12346, 274668, 12005168, 1018997864, 165091172592, 50502031367952}
	maxSize := 9 //Largest graphs we will search for.
	numberFound := make([]int, maxSize+1)
	m := 4
	for n := 0; n <= maxSize; n++ {
		for i := 0; i < m; i++ {
			iter := All(n, i, m)
			for iter.Next() {
				numberFound[n]++
			}
		}

	}
	t.Log(numberFound)
	if !ints.Equal(numberFound, truthData[:maxSize+1]) {
		t.Fail()
	}
}

func TestSaveAndLoad(t *testing.T) {
	//Initialise the preprune/prune functions.
	f := func(g *graph.DenseGraph) bool { return false }

	//Create a new iterator.
	iter := WithPruning(7, 0, 1, f, f)

	//Try saving the initial state
	initial := new(bytes.Buffer)
	iter.Save(initial)

	load := Load(initial, f, f)
	count := 0
	for load.Next() {
		count++
	}
	if count != 1044 {
		t.Fail()
	}

	for i := 0; i < 100; i++ {
		iter.Next()
	}

	//Try saving the state after 100 iterations.
	after100 := new(bytes.Buffer)
	iter.Save(after100)
	//Let's increment the state once and see if this throws it off.
	iter.Next()

	load = Load(after100, f, f)
	count = 0
	for load.Next() {
		count++
	}
	if count != 944 {
		t.Fail()
	}

	for iter.Next() {
	}

	//Try saving after the iterator is finished.
	end := new(bytes.Buffer)
	iter.Save(end)
	load = Load(end, f, f)
	count = 0
	for load.Next() {
		count++
	}
	if count != 0 {
		t.Fail()
	}

	//Try with the iteration split into two concurrent processes.
	iter0 := WithPruning(8, 0, 2, f, f)
	iter1 := WithPruning(8, 1, 2, f, f)
	count = 0
	for i := 0; i < 1000; i++ {
		iter0.Next()
		iter1.Next()
		count += 2
	}
	save0 := new(bytes.Buffer)
	save1 := new(bytes.Buffer)
	iter0.Save(save0)
	iter1.Save(save1)

	iter0 = Load(save0, f, f)
	iter1 = Load(save1, f, f)

	for iter0.Next() {
		count++
	}

	for iter1.Next() {
		count++
	}

	if count != 12346 {
		t.Fail()
	}
}

func BenchmarkAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iter := All(10, 0, 1)
		for iter.Next() {
		}
	}
}
