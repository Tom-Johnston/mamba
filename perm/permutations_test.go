package perm

import "testing"

func TestIterator(t *testing.T) {
	iter := NewIterator(Identity(4))
	correctResult := [][]int{[]int{0, 1, 2, 3}, []int{1, 0, 2, 3}, []int{2, 0, 1, 3}, []int{0, 2, 1, 3}, []int{1, 2, 0, 3}, []int{2, 1, 0, 3}, []int{3, 1, 0, 2}, []int{1, 3, 0, 2}, []int{0, 3, 1, 2}, []int{3, 0, 1, 2}, []int{1, 0, 3, 2}, []int{0, 1, 3, 2}, []int{0, 2, 3, 1}, []int{2, 0, 3, 1}, []int{3, 0, 2, 1}, []int{0, 3, 2, 1}, []int{2, 3, 0, 1}, []int{3, 2, 0, 1}, []int{3, 2, 1, 0}, []int{2, 3, 1, 0}, []int{1, 3, 2, 0}, []int{3, 1, 2, 0}, []int{2, 1, 3, 0}, []int{1, 2, 3, 0}}
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
		t.Log(v, w)
	}
}
