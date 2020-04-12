package itertools

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
