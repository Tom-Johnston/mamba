package itertools

import (
	"testing"

	"github.com/Tom-Johnston/mamba/ints"
)

func TestPermutations(t *testing.T) {
	iter := Permutations(4)
	correctResult := [][]int{{0, 1, 2, 3}, {1, 0, 2, 3}, {2, 0, 1, 3}, {0, 2, 1, 3}, {1, 2, 0, 3}, {2, 1, 0, 3}, {3, 1, 0, 2}, {1, 3, 0, 2}, {0, 3, 1, 2}, {3, 0, 1, 2}, {1, 0, 3, 2}, {0, 1, 3, 2}, {0, 2, 3, 1}, {2, 0, 3, 1}, {3, 0, 2, 1}, {0, 3, 2, 1}, {2, 3, 0, 1}, {3, 2, 0, 1}, {3, 2, 1, 0}, {2, 3, 1, 0}, {1, 3, 2, 0}, {3, 1, 2, 0}, {2, 1, 3, 0}, {1, 2, 3, 0}}
	i := -1
	for iter.Next() {
		i++
		v := iter.Value()
		w := correctResult[i]
		for j := range w {
			if w[j] != v[j] {
				t.Fail()
			}
		}
	}
	if i != 23 {
		t.Fail()
	}
}

func TestLexicographicPermutations(t *testing.T) {
	iter := LexicographicPermutations(4)
	correstResult := [][]int{{0, 1, 2, 3}, {0, 1, 3, 2}, {0, 2, 1, 3}, {0, 2, 3, 1}, {0, 3, 1, 2}, {0, 3, 2, 1}, {1, 0, 2, 3}, {1, 0, 3, 2}, {1, 2, 0, 3}, {1, 2, 3, 0}, {1, 3, 0, 2}, {1, 3, 2, 0}, {2, 0, 1, 3}, {2, 0, 3, 1}, {2, 1, 0, 3}, {2, 1, 3, 0}, {2, 3, 0, 1}, {2, 3, 1, 0}, {3, 0, 1, 2}, {3, 0, 2, 1}, {3, 1, 0, 2}, {3, 1, 2, 0}, {3, 2, 0, 1}, {3, 2, 1, 0}}
	i := -1
	for iter.Next() {
		i++
		v := iter.Value()
		w := correstResult[i]
		for j := range w {
			if w[j] != v[j] {
				t.Fail()
			}
		}
	}
	if i != 23 {
		t.Fail()
	}
}

func TestMultisetPermutations(t *testing.T) {
	iter := MultisetPermutations([]int{1, 2, 1})
	correstResult := [][]int{{0, 1, 1, 2}, {0, 1, 2, 1}, {0, 2, 1, 1}, {1, 0, 1, 2}, {1, 0, 2, 1}, {1, 1, 0, 2}, {1, 1, 2, 0}, {1, 2, 0, 1}, {1, 2, 1, 0}, {2, 0, 1, 1}, {2, 1, 0, 1}, {2, 1, 1, 0}}
	i := -1
	for iter.Next() {
		i++
		v := iter.Value()
		w := correstResult[i]
		for j := range w {
			if w[j] != v[j] {
				t.Fail()
			}
		}
	}
	if i != 11 {
		t.Fail()
	}

	iter = MultisetPermutations([]int{2, 2, 2})
	i = -1
	for iter.Next() {
		i++
	}
	if i != 89 {
		t.Fail()
	}
}

func TestTopologicalSorts(t *testing.T) {
	//This checks all 3x3 Young tableaux as in The Art of Computer programming.
	less := func(i, j int) bool {
		if i/3 == j/3 {
			return i < j
		}
		if i%3 == j%3 {
			return i < j
		}
		return false
	}
	iter := TopologicalSorts(9, less)
	correctResult := [][]int{{0, 1, 2, 3, 4, 5, 6, 7, 8},
		{0, 1, 2, 3, 4, 6, 5, 7, 8},
		{0, 1, 2, 3, 4, 7, 5, 6, 8},
		{0, 1, 2, 3, 5, 6, 4, 7, 8},
		{0, 1, 2, 3, 5, 7, 4, 6, 8},
		{0, 1, 3, 2, 4, 5, 6, 7, 8},
		{0, 1, 3, 2, 4, 6, 5, 7, 8},
		{0, 1, 3, 2, 4, 7, 5, 6, 8},
		{0, 1, 3, 2, 5, 6, 4, 7, 8},
		{0, 1, 3, 2, 5, 7, 4, 6, 8},
		{0, 1, 4, 2, 5, 6, 3, 7, 8},
		{0, 1, 4, 2, 5, 7, 3, 6, 8},
		{0, 1, 4, 2, 3, 5, 6, 7, 8},
		{0, 1, 4, 2, 3, 6, 5, 7, 8},
		{0, 1, 4, 2, 3, 7, 5, 6, 8},
		{0, 1, 5, 2, 3, 6, 4, 7, 8},
		{0, 1, 5, 2, 3, 7, 4, 6, 8},
		{0, 1, 6, 2, 3, 7, 4, 5, 8},
		{0, 1, 5, 2, 4, 6, 3, 7, 8},
		{0, 1, 5, 2, 4, 7, 3, 6, 8},
		{0, 1, 6, 2, 4, 7, 3, 5, 8},
		{0, 2, 3, 1, 4, 5, 6, 7, 8},
		{0, 2, 3, 1, 4, 6, 5, 7, 8},
		{0, 2, 3, 1, 4, 7, 5, 6, 8},
		{0, 2, 3, 1, 5, 6, 4, 7, 8},
		{0, 2, 3, 1, 5, 7, 4, 6, 8},
		{0, 2, 4, 1, 5, 6, 3, 7, 8},
		{0, 2, 4, 1, 5, 7, 3, 6, 8},
		{0, 3, 4, 1, 5, 6, 2, 7, 8},
		{0, 3, 4, 1, 5, 7, 2, 6, 8},
		{0, 2, 4, 1, 3, 5, 6, 7, 8},
		{0, 2, 4, 1, 3, 6, 5, 7, 8},
		{0, 2, 4, 1, 3, 7, 5, 6, 8},
		{0, 2, 5, 1, 3, 6, 4, 7, 8},
		{0, 2, 5, 1, 3, 7, 4, 6, 8},
		{0, 2, 6, 1, 3, 7, 4, 5, 8},
		{0, 2, 5, 1, 4, 6, 3, 7, 8},
		{0, 2, 5, 1, 4, 7, 3, 6, 8},
		{0, 2, 6, 1, 4, 7, 3, 5, 8},
		{0, 3, 5, 1, 4, 6, 2, 7, 8},
		{0, 3, 5, 1, 4, 7, 2, 6, 8},
		{0, 3, 6, 1, 4, 7, 2, 5, 8}}
	counter := 0
	for iter.Next() {
		if !ints.Equal(correctResult[counter], iter.InverseValue()) {
			t.Log(counter, iter.InverseValue(), correctResult[counter])
			t.Fail()
		}
		counter++
	}
}

func TestRestrictedPrefixPermutations(t *testing.T) {
	f := func(a []int) bool {
		if ints.Equal(a, []int{1}) || ints.Equal(a, []int{0, 2, 1}) || ints.Equal(a, []int{0, 3}) || ints.Equal(a, []int{2, 0, 3}) || ints.Equal(a, []int{3, 2, 0, 1}) {
			return false
		}
		return true
	}

	correctResult := [][]int{{0, 1, 2, 3}, {0, 1, 3, 2}, {0, 2, 3, 1}, {2, 0, 1, 3}, {2, 1, 0, 3}, {2, 1, 3, 0}, {2, 3, 0, 1}, {2, 3, 1, 0}, {3, 0, 1, 2}, {3, 0, 2, 1}, {3, 1, 0, 2}, {3, 1, 2, 0}, {3, 2, 1, 0}}

	iter := RestrictedPrefixPermutations(4, f)
	i := 0
	for iter.Next() {
		if !ints.Equal(iter.Value(), correctResult[i]) {
			t.FailNow()
		}
		i++
	}
	if i != len(correctResult) {
		t.Fail()
	}
}
