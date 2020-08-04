package itertools

//CombinationIterator iterates over all ways of choosing k elements from 0, ..., n-1 in lexicographic order.
type CombinationIterator struct {
	n    int
	k    int
	data []int
}

//Next moves the Iterator to the next combination, returning true if there is one and false is there isn't.
func (b *CombinationIterator) Next() bool {
	if b.k == 0 {
		b.k--
		return true
	}

	for i := b.k - 1; i >= 0; i-- {
		if b.data[i] < b.n+i-b.k {
			b.data[i]++
			for j := i + 1; j < b.k; j++ {
				b.data[j] = b.data[j-1] + 1
			}
			return true
		}
	}
	return false
}

//Value returns the current state of the iterator. This must not be modified.
func (b CombinationIterator) Value() []int {
	return b.data
}

//Combinations returns a new CombinationIterator which iterates over all subsets of k distinct elements from 0, ..., n-1 in lexicographic order.
func Combinations(n, k int) *CombinationIterator {
	data := make([]int, k)
	for i := 0; i < k; i++ {
		data[i] = i
	}

	if k > 0 {
		data[k-1]--
	}
	return &CombinationIterator{n: n, k: k, data: data}
}

//CombinationColexIterator iterates over all ways of chooding k elements from 0, ..., n-1 in colexicographic order.
//It can be initialised by calling CombinationsColex.
type CombinationColexIterator struct {
	n    int
	k    int
	j    int //Position to try to increase
	data []int
}

//CombinationsColex returns a new CombinationColexIterator which iterates over all subsets of k distinct elements from 0, ..., n-1 in colexicographic order.
func CombinationsColex(n, k int) *CombinationColexIterator {
	data := make([]int, k)
	for i := 0; i < k; i++ {
		data[i] = i
	}

	if k > 0 {
		data[k-1]--
	}
	return &CombinationColexIterator{n: n, k: k, j: k, data: data}
}

//Value returns the current subset. The output must not be modified.
func (b CombinationColexIterator) Value() []int {
	return b.data
}

//Next attempts to advance the iterator to the next subset, returning true if there is one and false if not.
func (b *CombinationColexIterator) Next() bool {
	if b.k <= 0 {
		b.k--
		return b.k == -1
	}

	if b.j >= b.k-1 {
		if b.data[b.k-1] == b.n-1 {
			return false
		}
		b.data[b.k-1]++
		b.j--
		return true
	}

	if b.j != -1 {
		b.data[b.j]++
		b.j--
		return true
	}

	for j := 0; j < b.k-1; j++ {
		if b.data[j] < b.data[j+1]-1 {
			b.data[j]++
			b.j = j - 1
			return true
		}
		b.data[j] = j
	}

	if b.data[b.k-1] == b.n-1 {
		return false
	}
	b.data[b.k-1]++
	b.j = b.k - 2
	return true
}
