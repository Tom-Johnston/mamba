package perm

//Permutation is the data structure for holding a permutation.
//A permutation p in on {0, 1, ..., n-1} is stored as a []int of length n where the ith element of the slice holds p(i).
type Permutation []int

//Identity returns the identity permutation which maps i to i.
func Identity(n int) []int {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return p
}

//Iterator is a struct containing the state of the iterator.
//It iterates over permutations according to Heap's algorithm.
type Iterator struct {
	n int
	i int
	c []int
	p []int
}

//NewIterator returns a new permuation iterator which iterates over all permutations of the elements in a.
//This doesn't modify a.
func NewIterator(a []int) *Iterator {
	copyOfA := make([]int, len(a))
	copy(copyOfA, a)
	p := Iterator{i: -1, n: len(a), c: make([]int, len(a)), p: copyOfA}
	return &p
}

//Value returns the current permutation.
//You must not modify the output of this function.
func (p *Iterator) Value() []int {
	return p.p
}

//Next moves the iterator to the next permutation, returning true if there is one and false if the previous permutation is the last one.
func (p *Iterator) Next() bool {
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
