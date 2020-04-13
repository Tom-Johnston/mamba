package itertools

import "testing"

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
