package itertools

import (
	"github.com/Tom-Johnston/mamba/ints"
)

//PermutationIterator is a struct containing the state of the iterator.
//It iterates over permutations according to Heap's algorithm.
type PermutationIterator struct {
	n int
	i int
	c []int
	p []int
}

//Permutations returns a new permuation iterator which iterates over all permutations of the elements {0, ..., n - 1}.
func Permutations(n int) *PermutationIterator {
	a := make([]int, n)
	for i := range a {
		a[i] = i
	}
	p := PermutationIterator{i: -1, n: len(a), c: make([]int, len(a)), p: a}
	return &p
}

//Value returns the current permutation.
//You must not modify the output of this function.
func (p *PermutationIterator) Value() []int {
	return p.p
}

//Next moves the iterator to the next permutation, returning true if there is one and false if the previous permutation is the last one.
func (p *PermutationIterator) Next() bool {
	if p.i == p.n {
		return false
	}

	if p.i == -1 {
		p.i++
		return true
	}

	for p.i < p.n {
		if p.c[p.i] < p.i {
			if p.i%2 == 0 {
				p.p[0], p.p[p.i] = p.p[p.i], p.p[0]
			} else {
				p.p[p.c[p.i]], p.p[p.i] = p.p[p.i], p.p[p.c[p.i]]
			}
			p.c[p.i]++
			p.i = 0
			return true
		}
		p.c[p.i] = 0
		p.i++
	}
	return false
}

//LexicographicPermutationIterator is a struct contraining the state of an iterator which iterates over all permutations of {0, 1, ..., n-1} in lexicographic order.
type LexicographicPermutationIterator struct {
	n     int
	a     []int
	first bool
}

//LexicographicPermutations returns a new iterator which iterates over all permutations of {0, ..., n-1} in lexicographic order.
//It is not safe to modify the output of the iterator.
func LexicographicPermutations(n int) *LexicographicPermutationIterator {
	a := make([]int, n)
	for i := range a {
		a[i] = i
	}
	return &LexicographicPermutationIterator{n: n, a: a, first: true}
}

//Value returns the current permutation.
//You must not modify the output of this function.
func (iter *LexicographicPermutationIterator) Value() []int {
	return iter.a
}

//Next moves the iterator to the next permutation, returning true if there is one and false if the previous permutation is the last one.
func (iter *LexicographicPermutationIterator) Next() bool {
	n := iter.n
	if n > 0 && iter.first {
		iter.first = false
		return true
	}

	if n > 1 && iter.a[n-2] < iter.a[n-1] {
		iter.a[n-2], iter.a[n-1] = iter.a[n-1], iter.a[n-2]
		return true
	}

	if n > 2 && iter.a[n-3] < iter.a[n-2] {
		if iter.a[n-3] < iter.a[n-1] {
			iter.a[n-3], iter.a[n-2], iter.a[n-1] = iter.a[n-1], iter.a[n-3], iter.a[n-2]
		} else {
			iter.a[n-3], iter.a[n-2], iter.a[n-1] = iter.a[n-2], iter.a[n-1], iter.a[n-3]
		}
		return true
	}

	for j := iter.n - 4; j >= 0; j-- {
		if iter.a[j] >= iter.a[j+1] {
			continue
		}
		if iter.a[j] < iter.a[n-1] {
			iter.a[j], iter.a[j+1], iter.a[n-1] = iter.a[n-1], iter.a[j], iter.a[j+1]
		} else {
			for l := n - 2; l > 0; l-- {
				if iter.a[j] >= iter.a[l] {
					continue
				}
				iter.a[j], iter.a[l] = iter.a[l], iter.a[j]
				iter.a[n-1], iter.a[j+1] = iter.a[j+1], iter.a[n-1]
				break
			}
		}
		k := j + 2
		l := n - 2
		for k < l {
			iter.a[k], iter.a[l] = iter.a[l], iter.a[k]
			k++
			l--
		}
		return true
	}
	return false
}

//MultisetPermutationIterator is a struct containing the state of an iterator which iterates over all permutations of some multiset from {0, 1, ..., n-1} in lexicographic order.
type MultisetPermutationIterator struct {
	lexIter *LexicographicPermutationIterator
}

//MultisetPermutations returns an iterator which iterates over all permutations of some multiset from {0, 1, ..., n-1} in lexicographic order. The number i appears freq[i] times in the multiset.
//It is not safe to modify the output of the iterator.
func MultisetPermutations(freq []int) *MultisetPermutationIterator {
	n := ints.Sum(freq)
	a := make([]int, 0, n)
	for i := range freq {
		for j := 0; j < freq[i]; j++ {
			a = append(a, i)
		}
	}
	return &MultisetPermutationIterator{lexIter: &LexicographicPermutationIterator{n: n, a: a, first: true}}
}

//Value returns the current permutation.
//You must not modify the output of this function.
func (iter *MultisetPermutationIterator) Value() []int {
	return iter.lexIter.Value()
}

//Next moves the iterator to the next permutation, returning true if there is one and false if the previous permutation is the last one.
func (iter *MultisetPermutationIterator) Next() bool {
	return iter.lexIter.Next()
}

//Topological sorts as in Algorithm V of The Art of Computer Programming Volume 4a Section 7.2.1.2.

//TopologicalSortIterator iterates over all topological sorts of {0, 1, ..., n-1} respecting some partial order of the total order 0 < 1 < ... < n -1.
type TopologicalSortIterator struct {
	less     func(i, j int) bool
	state    []int
	invState []int
	n        int
	first    bool
}

//TopologicalSorts returns an iterator which iterates over all topological sorts of {0, 1, ... , n-1} according to the partial order less. If less(i,j) == true, then this only iterates over permutations where i appears before j.
//less must be a sub-order of the total order 0 < 1 < ... n -1 i.e. for inputs i < j the function less should return true if the condition i < j should be imposed and false otherwise. If i >= j, the function should return false.
//The function less doesn't need to be transitive as if i < j < k, then less(i,k) will never be called.
//If less(i,j) == true, then the inverse permutation is such that the value in position i is smaller than the value in position j. The function iter.InverseValue() returns the inverse permutation.
func TopologicalSorts(n int, less func(i, j int) bool) *TopologicalSortIterator {
	state := make([]int, n)
	for i := range state {
		state[i] = i
	}
	invState := make([]int, n)
	for i := range invState {
		invState[i] = i
	}

	return &TopologicalSortIterator{less: less, state: state, invState: invState, n: n, first: true}
}

//Value returns the current permutation.
//You must not modify the value.
func (iter *TopologicalSortIterator) Value() []int {
	return iter.state
}

//InverseValue returns the inverse the to current permutation.
//You must not modify the value.
func (iter *TopologicalSortIterator) InverseValue() []int {
	return iter.invState
}

//Next attempts to advance the iterator to the next permutation, returning true if there is one and false otherwise.
func (iter *TopologicalSortIterator) Next() bool {
	if iter.first {
		iter.first = false
		return true
	}

	n := iter.n
	for k := n - 1; k >= 0; k-- {
		j := iter.invState[k]
		if j > 0 {
			l := iter.state[j-1]
			//TODO: Do we always have l < k here?
			if !iter.less(l, k) {
				iter.state[j-1] = k
				iter.state[j] = l
				iter.invState[k] = j - 1
				iter.invState[l] = j
				return true
			}
		}

		for j < k {
			l := iter.state[j+1]
			iter.state[j] = l
			iter.invState[l] = j
			j++
		}
		iter.state[k] = k
		iter.invState[k] = k
	}
	return false
}
