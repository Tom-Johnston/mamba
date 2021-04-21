package itertools

import (
	"sort"
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

func TestPermutationsByPattern(t *testing.T) {
	//flatten replaces src by the unique permutation which is order-isomorphic to src.
	//If src[i] is the jth smallest entry amongst all the elements in src, then src[i] is set to j.
	//buf is storage space that may be used in the computation. buf must have capacity at least len(src).
	flatten := func(src, buf []int) {
		//This is a very basic implementation and could undoubtedly be quicker.
		n := len(src)
		buf = buf[:n]
		copy(buf, src)
		ints.Sort(buf)

		for i, v := range src {
			u := sort.SearchInts(buf, v)
			src[i] = u
		}
	}

	//endsWithSquare checks if the permutation contains a square which ends at the last entry of a.
	//buf is storage space which may be used in the computation. buf must have capacity at least 3*len(a).
	endsWithSquare := func(a []int, buf []int) bool {
		n := len(a)

		//Check for squares of length 2.
		if n >= 4 {
			if a[n-2] < a[n-1] && a[n-4] < a[n-3] {
				return true
			}
			if a[n-2] > a[n-1] && a[n-4] > a[n-3] {
				return true
			}
		}

		//Now check for squares of larger lengths.
		//Since we know the permutation has no squares of length 2, the length of any square must be a multiple of 4.
		for i := 4; i <= n/2; i += 4 {
			patt1 := buf[:i]
			patt2 := buf[n : n+i]
			copy(patt1, a[n-i:n])
			copy(patt2, a[n-2*i:n-i])
			flatten(patt1, buf[2*n:])
			flatten(patt2, buf[2*n:])
			if ints.Equal(patt1, patt2) {
				return true
			}
		}
		return false
	}

	truthData := []int{1, 1, 2, 6, 12, 34, 104, 406, 1112, 3980, 15216, 68034, 312048, 1625968, 8771376, 53270068, 319218912, 2135312542, 14420106264}
	for i := 0; i < 10; i++ {
		v := truthData[i]
		buf := make([]int, 3*i)
		f := func(a []int) bool { return !endsWithSquare(a, buf) }
		iter := PermutationsByPattern(i, f)
		count := 0
		for iter.Next() {
			count++
		}
		if count != v {
			t.Log(count)
			t.Fail()
		}
	}
}
