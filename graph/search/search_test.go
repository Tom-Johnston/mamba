package search

import (
	"testing"

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

func BenchmarkAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iter := All(10, 0, 1)
		for iter.Next() {
		}
	}
}
